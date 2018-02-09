// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Flexvolume driver that is invoked by kubelet when a pod installs a flexvolume drive
// of type nodeagent/uds
// This driver communicates to the nodeagent using either
//   * (Default) writing credentials of workloads to a file or
//   * gRPC message defined at protos/nodeagementmgmt.proto,
// to shares the properties of the pod with nodeagent.
//
package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	nagent "istio.io/istio/security/cmd/node_agent_k8s/nodeagentmgmt"
	pb "istio.io/istio/security/proto"
)

// Response is the output of Flex volume driver to the kubelet.
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	// Is attached resp.
	Attached bool `json:"attached,omitempty"`
	// Dev mount resp.
	Device string `json:"device,omitempty"`
	// Volumen name resp.
	VolumeName string `json:"volumename,omitempty"`
}

// Response to the 'init' command.
// We want to explicitly set and send Attach: false
// that is why it is separated from the Response struct.
type InitResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	// Capability resp.
	Attach bool `json:"attach"`
}

// ConfigurationOptions may be used to setup the driver.
// These are optional and most users will not depened on them and will instead use the defaults.
type ConfigurationOptions struct {
	// Version of the Kubernetes cluster on which the driver is running.
	K8sVersion string `json:"k8s_version,omitempty"`
	// Location on the node's filesystem where the driver will host the
	// per workload directory and the credentials for the workload.
	// Default: /var/run/nodeagent
	NodeAgentManagementHomeDir string `json:"nodeagent_management_home,omitempty"`
	// Relative location to NodeAgentManagementHomeDir where per workload directory
	// will be created.
	// Default: /mount
	// For example: /mount here implies /var/run/nodeagent/mount/ directory
	// on the node.
	NodeAgentWorkloadHomeDir string `json:"nodeagent_workload_home,omitempty"`
	// Relative location to NodeAgentManagementHomeDir where per workload credential
	// files will be created.
	// Default: /creds
	// For example: /creds here implies /var/run/nodeagent/creds/ directory
	NodeAgentCredentialsHomeDir string `json:"nodeagent_credentials_home,omitempty"`
	// Whether node agent will be notified of workload using gRPC or by creating files in
	// NodeAgentCredentialsHomeDir.
	// Default: false
	UseGrpc bool `json:"use_grpc,omitempty"`
	// If UseGrpc is true then host the UDS socket here.
	// This is relative to NodeAgentManagementHomeDir
	// Default: mgmt.sock
	// For example: /mgmt/mgmt.sock implies /var/run/nodeagement/mgmt/mgmt.sock
	NodeAgentManagementApi string `json:"nodeagent_management_api,omitempty"`
	// Log level for loggint to node syslog. Options: INFO|WARNING
	// Default: WARNING
	LogLevel string `json:"log_level,omitempty"`
}

// FlexVolumeInputs is the structure used by kubelet to notify driver
// volume mounts/unmounts.
type FlexVolumeInputs struct {
	Uid            string `json:"kubernetes.io/pod.uid"`
	Name           string `json:"kubernetes.io/pod.name"`
	Namespace      string `json:"kubernetes.io/pod.namespace"`
	ServiceAccount string `json:"kubernetes.io/serviceAccount.name"`
}

const (
	VER_K8S        string = "1.8"
	CONFIG_FILE    string = "/etc/flexvolume/nodeagent.json"
	NODEAGENT_HOME string = "/var/run/nodeagent"
	MOUNT_DIR      string = "/mount"
	CREDS_DIR      string = "/creds"
	MGMT_SOCK      string = "/mgmt.sock"
	LOG_LEVEL_WARN string = "WARNING"
)

var (
	// configuration is the active configuration that is being used by the driver.
	configuration *ConfigurationOptions
	// logWriter is used to notify syslog of the functionality of the driver.
	logWriter *syslog.Writer
	// defaultConfiguration is the default configuration for the driver.
	defaultConfiguration ConfigurationOptions = ConfigurationOptions{
		K8sVersion:                  VER_K8S,
		NodeAgentManagementHomeDir:  NODEAGENT_HOME,
		NodeAgentWorkloadHomeDir:    MOUNT_DIR,
		NodeAgentCredentialsHomeDir: CREDS_DIR,
		UseGrpc:                     false,
		NodeAgentManagementApi:      MGMT_SOCK,
		LogLevel:                    LOG_LEVEL_WARN,
	}
	configFile string = CONFIG_FILE
	// setup as var's so that we can test.
	getExecCmd = exec.Command
)

// initCommand handles the init command for the driver.
func InitCommand() error {
	if configuration.K8sVersion == "1.8" {
		resp, err := json.Marshal(&InitResponse{Status: "Success", Message: "Init ok.", Attach: false})
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	}
	return genericSuccess("init", "", "Init ok.")
}

// checkValidMountOpts checks if there are sufficient inputs to
// call node agent.
func checkValidMountOpts(opts string) (*pb.WorkloadInfo, bool) {
	ninputs := FlexVolumeInputs{}
	err := json.Unmarshal([]byte(opts), &ninputs)
	if err != nil {
		return nil, false
	}

	wlInfo := pb.WorkloadInfo{
		Attrs: &pb.WorkloadInfo_WorkloadAttributes{
			Uid:            ninputs.Uid,
			Workload:       ninputs.Name,
			Namespace:      ninputs.Namespace,
			Serviceaccount: ninputs.ServiceAccount,
		},
		Workloadpath: ninputs.Uid,
	}
	return &wlInfo, true
}

// doMount handles a new workload mounting the flex volume drive. It will:
// * mount a tmpfs at the destinationDir(ectory) of the workload created by the kubelet.
// * create a sub-directory ('nodeagent') there
// * do a bind mount of the nodeagent's directory on the node to the destinationDir/nodeagent.
func doMount(destinationDir string, ninputs *pb.WorkloadInfo) error {
	newDir := filepath.Join(configuration.NodeAgentWorkloadHomeDir, ninputs.Workloadpath)
	err := os.MkdirAll(newDir, 0777)
	if err != nil {
		return err
	}

	// Not really needed but attempt to workaround:
	// https://github.com/kubernetes/kubernetes/blob/61ac9d46382884a8bd9e228da22bca5817f6d226/pkg/util/mount/mount_linux.go
	cmdMount := getExecCmd("/bin/mount", "-t", "tmpfs", "-o", "size=8K", "tmpfs", destinationDir)
	err = cmdMount.Run()
	if err != nil {
		os.RemoveAll(newDir)
		return err
	}

	newDestinationDir := filepath.Join(destinationDir, "nodeagent")
	err = os.MkdirAll(newDestinationDir, 0777)
	if err != nil {
		cmd := getExecCmd("/bin/unmount", destinationDir)
		cmd.Run()
		os.RemoveAll(newDir)
		return err
	}

	// Do a bind mount
	cmd := getExecCmd("/bin/mount", "--bind", newDir, newDestinationDir)
	err = cmd.Run()
	if err != nil {
		cmd = getExecCmd("/bin/umount", destinationDir)
		cmd.Run()
		os.RemoveAll(newDir)
		return err
	}

	return nil
}

// doUnmount will unmount the directory
func doUnmount(dir string) error {
	cmd := getExecCmd("/bin/umount", dir)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Mount handles the mount command to the driver.
func Mount(dir, opts string) error {
	inp := strings.Join([]string{dir, opts}, "|")

	ninputs, s := checkValidMountOpts(opts)
	if s == false {
		return failure("mount", inp, "Incomplete inputs")
	}

	if err := doMount(dir, ninputs); err != nil {
		sErr := "Failure to mount: " + err.Error()
		return failure("mount", inp, sErr)
	}

	if configuration.UseGrpc == true {
		if err := sendWorkloadAdded(ninputs); err != nil {
			sErr := "Failure to notify nodeagent: " + err.Error()
			return failure("mount", inp, sErr)
		}
	} else {
		if err := addCredentialFile(ninputs); err != nil {
			sErr := "Failure to create credentials: " + err.Error()
			return failure("mount", inp, sErr)
		}
	}

	return genericSuccess("mount", inp, "Mount ok.")
}

// unmount handles the unmount command to the driver.
func Unmount(dir string) error {
	var emsgs []string
	// Stop the listener.
	// /var/lib/kubelet/pods/20154c76-bf4e-11e7-8a7e-080027631ab3/volumes/nodeagent~uds/test-volume/
	// /var/lib/kubelet/pods/2dc75e9a-cbec-11e7-b158-0800270da466/volumes/nodeagent~uds/test-volume
	comps := strings.Split(dir, "/")
	if len(comps) < 6 {
		sErr := fmt.Sprintf("Failure to notify nodeagent dir %v", dir)
		return failure("unmount", dir, sErr)
	}

	uid := comps[5]
	// TBD: Check if uid is the correct format.
	naInp := &pb.WorkloadInfo{
		Attrs:        &pb.WorkloadInfo_WorkloadAttributes{Uid: uid},
		Workloadpath: uid,
	}
	if configuration.UseGrpc == true {
		if err := sendWorkloadDeleted(naInp); err != nil {
			sErr := "Failure to notify nodeagent: " + err.Error()
			return failure("unmount", dir, sErr)
		}
	} else {
		if err := removeCredentialFile(naInp); err != nil {
			// Go ahead and finish the unmount; no need to hold up kubelet.
			emsgs = append(emsgs, "Failure to delete credentials file: "+err.Error())
		}
	}

	// unmount the bind mount
	doUnmount(filepath.Join(dir, "nodeagent"))
	// unmount the tmpfs
	doUnmount(dir)
	// delete the directory that was created.
	delDir := filepath.Join(configuration.NodeAgentWorkloadHomeDir, uid)
	err := os.Remove(delDir)
	if err != nil {
		emsgs = append(emsgs, fmt.Sprintf("unmount del failure %s: %s", delDir, err.Error()))
		// go head and return ok.
	}

	if len(emsgs) == 0 {
		emsgs = append(emsgs, "Unmount Ok")
	}

	return genericSuccess("unmount", dir, strings.Join(emsgs, ","))
}

// printAndLog is used to print to stdout and to the syslog.
func printAndLog(caller, inp, s string) {
	fmt.Println(s)
	logToSys(caller, inp, s)
}

// genericSuccess is to print a success response to the kubelet.
func genericSuccess(caller, inp, msg string) error {
	resp, err := json.Marshal(&Response{Status: "Success", Message: msg})
	if err != nil {
		return err
	}

	printAndLog(caller, inp, string(resp))
	return nil
}

// failure is to print a failure response to the kubelet.
func failure(caller, inp, msg string) error {
	resp, err := json.Marshal(&Response{Status: "Failure", Message: msg})
	if err != nil {
		return errors.New(msg)
	}

	printAndLog(caller, inp, string(resp))
	return errors.New(msg)
}

// GenericUnsupported is to print a un-supported response to the kubelet.
func GenericUnsupported(caller, inp, msg string) error {
	resp, err := json.Marshal(&Response{Status: "Not supported", Message: msg})
	if err != nil {
		return err
	}

	printAndLog(caller, inp, string(resp))
	return nil
}

// logToSys is to write to syslog.
func logToSys(caller, inp, opts string) {
	if logWriter == nil {
		return
	}

	opt := strings.Join([]string{caller, inp, opts}, "|")
	if configuration.LogLevel == LOG_LEVEL_WARN {
		logWriter.Warning(opt)
	} else {
		logWriter.Info(opt)
	}
}

// sendWorkloadAdded is used to notify node agent of a addition of a workload with the flex-volume volume mounted.
func sendWorkloadAdded(ninputs *pb.WorkloadInfo) error {
	client := nagent.ClientUds(configuration.NodeAgentManagementApi)
	if client == nil {
		return errors.New("Failed to create Nodeagent client.")
	}

	_, err := client.WorkloadAdded(ninputs)
	if err != nil {
		return err
	}

	client.Close()

	return nil
}

// sendWorkloadDeleted is used to notify node agent of deletion of a workload with the flex-volume volume mounted.
func sendWorkloadDeleted(ninputs *pb.WorkloadInfo) error {
	client := nagent.ClientUds(configuration.NodeAgentManagementApi)
	if client == nil {
		return errors.New("Failed to create Nodeagent client.")
	}

	_, err := client.WorkloadDeleted(ninputs)
	if err != nil {
		return err
	}

	client.Close()
	return nil
}

func getCredFile(uid string) string {
	return uid + ".json"
}

// addCredentialFile is used to create a credential file when a workload with the flex-volume volume mounted is created.
func addCredentialFile(ninputs *pb.WorkloadInfo) error {
	//Make the directory and then write the ninputs as json to it.
	var err error
	err = os.MkdirAll(configuration.NodeAgentCredentialsHomeDir, 755)
	if err != nil {
		return err
	}

	var attrs []byte
	attrs, err = json.Marshal(ninputs.Attrs)
	if err != nil {
		return err
	}

	credsFileTmp := filepath.Join(configuration.NodeAgentManagementHomeDir, getCredFile(ninputs.Attrs.Uid))
	err = ioutil.WriteFile(credsFileTmp, attrs, 0644)

	// Move it to the right location now.
	credsFile := filepath.Join(configuration.NodeAgentCredentialsHomeDir, getCredFile(ninputs.Attrs.Uid))
	return os.Rename(credsFileTmp, credsFile)
}

// removeCredentialFile is used to delete a credential file when a workload with the flex-volume volume mounted is deleted.
func removeCredentialFile(ninputs *pb.WorkloadInfo) error {
	credsFile := filepath.Join(configuration.NodeAgentCredentialsHomeDir, getCredFile(ninputs.Attrs.Uid))
	err := os.Remove(credsFile)
	return err
}

//mkAbsolutePaths converts all the configuration paths to be absolute.
func mkAbsolutePaths(config *ConfigurationOptions) {
	config.NodeAgentWorkloadHomeDir = filepath.Join(config.NodeAgentManagementHomeDir, config.NodeAgentWorkloadHomeDir)
	config.NodeAgentCredentialsHomeDir = filepath.Join(config.NodeAgentManagementHomeDir, config.NodeAgentCredentialsHomeDir)
	config.NodeAgentManagementApi = filepath.Join(config.NodeAgentManagementHomeDir, config.NodeAgentManagementApi)
}

// If available read the configuration file and initialize the configuration options
// of the driver.
func InitConfiguration() {
	configuration = &defaultConfiguration
	mkAbsolutePaths(configuration)

	if _, err := os.Stat(configFile); err != nil {
		// Return quietly
		return
	}

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		logWriter.Warning(fmt.Sprintf("Not able to read %s: %s\n", CONFIG_FILE, err.Error()))
		return
	}

	var config ConfigurationOptions
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		logWriter.Warning(fmt.Sprintf("Not able to parse %s: %s\n", CONFIG_FILE, err.Error()))
		return
	}

	//fill in if missing configurations
	if len(config.NodeAgentManagementHomeDir) == 0 {
		config.NodeAgentManagementHomeDir = NODEAGENT_HOME
	}

	if len(config.NodeAgentWorkloadHomeDir) == 0 {
		config.NodeAgentWorkloadHomeDir = MOUNT_DIR
	}

	if len(config.NodeAgentCredentialsHomeDir) == 0 {
		config.NodeAgentCredentialsHomeDir = CREDS_DIR
	}

	if len(config.NodeAgentManagementApi) == 0 {
		config.NodeAgentManagementApi = MGMT_SOCK
	}

	if len(config.LogLevel) == 0 {
		config.LogLevel = LOG_LEVEL_WARN
	}

	if len(config.K8sVersion) == 0 {
		config.K8sVersion = VER_K8S
	}

	configuration = &config
	mkAbsolutePaths(configuration)
}

func logLevel(level string) syslog.Priority {
	switch level {
	case LOG_LEVEL_WARN:
		return syslog.LOG_WARNING
	}
	return syslog.LOG_INFO
}

func InitLog(tag string) (*syslog.Writer, error) {
	logWriter, err := syslog.New(logLevel(configuration.LogLevel)|syslog.LOG_DAEMON, tag)
	return logWriter, err
}
