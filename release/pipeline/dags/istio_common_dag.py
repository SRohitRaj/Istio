"""Airfow DAG and helpers used in one or more istio release pipeline."""
"""Copyright 2017 Istio Authors. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import datetime
import logging
import time
import collections

from airflow import DAG
from airflow.models import Variable
from airflow.operators.bash_operator import BashOperator
from airflow.operators.dummy_operator import DummyOperator
from airflow.operators.python_operator import BranchPythonOperator
from airflow.operators.python_operator import PythonOperator

import environment_config
from gcs_copy_operator import GoogleCloudStorageCopyOperator

default_args = {
    'owner': 'rkrishnap',
    'depends_on_past': False,
    # This is the date to when the airflow pipeline thinks the run started
    'start_date': datetime.datetime.now() - datetime.timedelta(days=1, minutes=15),
    'email': environment_config.EMAIL_LIST,
    'email_on_failure': True, 
    'email_on_retry': False,
    'retries': 1,
    'retry_delay': datetime.timedelta(minutes=5),
}

CommonTasks = collections.OrderedDict([
	('generate_workflow_args',          ''),
	('get_git_commit',                  ''),
	('run_cloud_builder',               ''),
	('run_release_qualification_tests', ''),
	('copy_files_for_release',          '')])
#        'mark_daily_complete'       used by   daily release
#        'mark_monthly_complete'     used by monthly release
#        'github_and_docker_release' used by monthly release
#        'github_tag_repos'          used by monthly release

def GetSettingPython(ti, setting):
  """Get setting form the generate_flow_args task.

  Args:
    ti: (task_instance) This is provided by the environment
    setting: (string) The name of the setting.
  Returns:
    The item saved in xcom.
  """
  return ti.xcom_pull(task_ids='generate_workflow_args')[setting]


def GetSettingTemplate(setting):
  """Create the template that will resolve to a setting from xcom.

  Args:
    setting: (string) The name of the setting.
  Returns:
    A templated string that resolves to a setting from xcom.
  """
  return ('{{ task_instance.xcom_pull(task_ids="generate_workflow_args"'
          ').%s }}') % (
              setting)


def AirflowGetVariableOrBaseCase(var, base):
  try:
    return Variable.get(var)
  except KeyError:
    return base

def MakeCommonDag(dag_args_func, name='istio_daily_flow_test',
                  schedule_interval='15 9 * * *'):
  """Creates the shared part of the daily/release dags."""
  common_dag = DAG(
      name,
      catchup=False,
      default_args=default_args,
      schedule_interval=schedule_interval,
  )
  tasksOD = collections.OrderedDict(CommonTasks)

  generate_flow_args = PythonOperator(
      task_id='generate_workflow_args',
      python_callable=dag_args_func,
      provide_context=True,
      dag=common_dag,
  )
  tasksOD['generate_workflow_args'] = generate_flow_args

  get_git_commit_cmd = """
    {% set settings = task_instance.xcom_pull(task_ids='generate_workflow_args') %}
    git config --global user.name "TestRunnerBot"
    git config --global user.email "testrunner@istio.io"
    git clone {{ settings.MFEST_URL }} green-builds || exit 2

    pushd green-builds
    git checkout {{ settings.BRANCH }}
    git checkout {{ settings.MFEST_COMMIT }} || exit 3
    ISTIO_SHA=`grep {{ settings.GITHUB_ORG }}/{{ settings.GITHUB_REPO }} {{ settings.MFEST_FILE }} | cut -f 6 -d \\"` || exit 4
    API_SHA=`  grep {{ settings.GITHUB_ORG }}/api                        {{ settings.MFEST_FILE }} | cut -f 6 -d \\"` || exit 5
    PROXY_SHA=`grep {{ settings.GITHUB_ORG }}/proxy                      {{ settings.MFEST_FILE }} | cut -f 6 -d \\"` || exit 6
    if [ -z ${ISTIO_SHA} ] || [ -z ${API_SHA} ] || [ -z ${PROXY_SHA} ]; then
      echo "ISTIO_SHA:$ISTIO_SHA API_SHA:$API_SHA PROXY_SHA:$PROXY_SHA some shas not found"
      exit 7
    fi
    popd #green-builds

    git clone {{ settings.ISTIO_REPO }} istio-code -b {{ settings.BRANCH }}
    pushd istio-code/release
    ISTIO_HEAD_SHA=`git rev-parse HEAD`
    git checkout ${ISTIO_SHA} || exit 8
    
    TS_SHA=` git show -s --format=%ct ${ISTIO_SHA}`
    TS_HEAD=`git show -s --format=%ct ${ISTIO_HEAD_SHA}`
    DIFF_SEC=$((TS_HEAD - TS_SHA))
    DIFF_DAYS=$(($DIFF_SEC/86400))
    if [ "{{ settings.CHECK_GREEN_SHA_AGE }}" = "true" ] && [ "$DIFF_DAYS" -gt "2" ]; then
       echo ERROR: ${ISTIO_SHA} is $DIFF_DAYS days older than head of branch {{ settings.BRANCH }}
       exit 9
    fi
    popd #istio-code/release

    if [ "{{ settings.VERIFY_CONSISTENCY }}" = "true" ]; then
      PROXY_REPO=`dirname {{ settings.ISTIO_REPO }}`/proxy
      echo $PROXY_REPO
      git clone $PROXY_REPO proxy-code -b {{ settings.BRANCH }}
      pushd proxy-code
      PROXY_HEAD_SHA=`git rev-parse HEAD`
      PROXY_HEAD_API_SHA=`grep ISTIO_API istio.deps  -A 4 | grep lastStableSHA | cut -f 4 -d '"'`
      popd
      if [ "$PROXY_HEAD_SHA" != "$PROXY_SHA" ]; then
        echo "inconsistent shas     PROXY_HEAD_SHA     $PROXY_HEAD_SHA != $PROXY_SHA PROXY_SHA" 1>&2
        exit 10
      fi
      if [ "$PROXY_HEAD_API_SHA" != "$API_SHA" ]; then
        echo "inconsistent shas PROXY_HEAD_API_SHA $PROXY_HEAD_API_SHA != $API_SHA   API_SHA"   1>&2 
        exit 11
      fi
      if [ "$ISTIO_HEAD_SHA" != "$ISTIO_SHA" ]; then
        echo "inconsistent shas     ISTIO_HEAD_SHA     $ISTIO_HEAD_SHA != $ISTIO_SHA ISTIO_SHA" 1>&2 
        exit 12
      fi
    fi

    pushd istio-code/release
    gsutil cp *.sh   gs://{{ settings.GCS_RELEASE_TOOLS_PATH }}/data/release/
    gsutil cp *.json gs://{{ settings.GCS_RELEASE_TOOLS_PATH }}/data/release/
    popd #istio-code/release

    pushd green-builds
    git rev-parse HEAD
    """

  get_git_commit = BashOperator(
      task_id='get_git_commit',
      bash_command=get_git_commit_cmd,
      xcom_push=True,
      dag=common_dag)
  tasksOD['get_git_commit'] = get_git_commit

  build_template = """
    {% set settings = task_instance.xcom_pull(task_ids='generate_workflow_args') %}
    {% set m_commit = task_instance.xcom_pull(task_ids='get_git_commit') %}
    gsutil cp gs://{{ settings.GCS_RELEASE_TOOLS_PATH }}/data/release/*.json .
    gsutil cp gs://{{ settings.GCS_RELEASE_TOOLS_PATH }}/data/release/*.sh .
    chmod u+x *
    ./start_gcb_build.sh -w -p {{ settings.PROJECT_ID \
    }} -r {{ settings.GCR_STAGING_DEST }} -s {{ settings.GCS_BUILD_PATH }} \
    -v "{{ settings.VERSION }}" \
    -u "{{ settings.MFEST_URL }}" \
    -t "{{ m_commit }}" -m "{{ settings.MFEST_FILE }}" \
    -a {{ settings.SVC_ACCT }}
    """
  # NOTE: if you add commands to build_template after start_gcb_build.sh then take care to preserve its return value

  build = BashOperator(
      task_id='run_cloud_builder', bash_command=build_template, dag=common_dag)
  tasksOD['run_cloud_builder'] = build

  test_command = """
    cp /home/airflow/gcs/data/githubctl ./githubctl
    chmod u+x ./githubctl
    {% set settings = task_instance.xcom_pull(task_ids='generate_workflow_args') %}
    git config --global user.name "TestRunnerBot"
    git config --global user.email "testrunner@istio.io"
    ls -l    ./githubctl
    ./githubctl \
    --token_file="{{ settings.TOKEN_FILE }}" \
    --op=dailyRelQual \
    --hub=gcr.io/{{ settings.GCR_STAGING_DEST }} \
    --gcs_path="{{ settings.GCS_BUILD_PATH }}" \
    --tag="{{ settings.VERSION }}" \
    --base_branch="{{ settings.BRANCH }}"
    """

  run_release_qualification_tests = BashOperator(
      task_id='run_release_qualification_tests',
      bash_command=test_command,
      retries=0,
      dag=common_dag)
  tasksOD['run_release_qualification_tests'] = run_release_qualification_tests

  copy_files = GoogleCloudStorageCopyOperator(
      task_id='copy_files_for_release',
      source_bucket=GetSettingTemplate('GCS_BUILD_BUCKET'),
      source_object=GetSettingTemplate('GCS_STAGING_PATH'),
      destination_bucket=GetSettingTemplate('GCS_STAGING_BUCKET'),
      dag=common_dag,
  )
  tasksOD['copy_files_for_release'] = copy_files

  return common_dag, tasksOD


def ReportDailySuccessful(task_instance, **kwargs):
  """Set this release as the candidate if it is the latest."""
  date = kwargs['execution_date']
  branch = GetSettingPython(task_instance, 'BRANCH')
  latest_run = float(AirflowGetVariableOrBaseCase(branch+'latest_daily_timestamp', 0))

  timestamp = time.mktime(date.timetuple())
  logging.info("Current run's timestamp: %s \n"
               "latest_daily's timestamp: %s", timestamp, latest_run)
  if timestamp >= latest_run:
    run_sha = task_instance.xcom_pull(task_ids='get_git_commit')
    latest_version = GetSettingPython(task_instance, 'VERSION')

    Variable.set(branch+'_latest_sha', run_sha)
    Variable.set(branch+'_latest_daily', latest_version)
    Variable.set(branch+'_latest_daily_timestamp', timestamp)

    logging.info('%s_latest_sha set to %s', branch, run_sha)
    logging.info('setting latest green daily of %s branch to: %s', branch, run_sha)
    return 'tag_daily_gcr'
  return 'skip_tag_daily_gcr'


def MakeMarkComplete(dag):
  """Make the final sequence of the daily graph."""
  mark_complete = BranchPythonOperator(
      task_id='mark_complete',
      python_callable=ReportDailySuccessful,
      provide_context=True,
      dag=dag,
  )

  gcr_tag_success = r"""
{% set settings = task_instance.xcom_pull(task_ids='generate_workflow_args') %}
set -x
pwd; ls

gsutil ls gs://{{ settings.GCS_FULL_STAGING_PATH }}/docker/           > docker_tars.txt
  cat docker_tars.txt |   grep -Eo "docker\/(([a-z]|-)*).tar.gz" | \
                          sed -E "s/docker\/(([a-z]|-)*).tar.gz/\1/g" > docker_images.txt

  gcloud auth configure-docker  -q
  cat docker_images.txt | \
  while read -r docker_image;do
    pull_source="gcr.io/{{ settings.GCR_STAGING_DEST }}/${docker_image}:{{ settings.VERSION }}"
    push_dest="  gcr.io/{{ settings.GCR_STAGING_DEST }}/${docker_image}:latest_{{ settings.BRANCH }}";
    docker pull $pull_source
    docker tag  $pull_source $push_dest
    docker push $push_dest
  done

cat docker_tars.txt docker_images.txt
rm  docker_tars.txt docker_images.txt
"""

  tag_daily_grc = BashOperator(
      task_id='tag_daily_gcr',
      bash_command=gcr_tag_success,
      dag=dag,
  )
  # skip_grc = DummyOperator(
  #     task_id='skip_tag_daily_gcr',
  #     dag=dag,
  # )
  # end = DummyOperator(
  #     task_id='end',
  #     dag=dag,
  #     trigger_rule="one_success",
  # )
  mark_complete >> tag_daily_grc
  # mark_complete >> skip_grc >> end
  return mark_complete



def DailyPipeline(branch):
  def DailyGenerateTestArgs(**kwargs):
    daily_params = dict(
      GCS_DAILY_PATH='daily-build/{version}',
      MFEST_COMMIT='{branch}@{{{timestamp}}}',
      VERIFY_CONSISTENCY='false', 
      VERSION='{branch}-{date_string}')

    """Loads the configuration that will be used for this Iteration."""
    conf = kwargs['dag_run'].conf
    if conf is None:
      conf = dict()

    # If variables are overriden then we should use it otherwise we use it's
    # default value.
    date = datetime.datetime.now()
    date_string = date.strftime('%Y%m%d-%H-%M')

    version = conf.get('VERSION')
    if version is None:
      version = daily_params['VERSION'].format(
        branch=branch, date_string=date_string)

    gcs_path = conf.get('GCS_DAILY_PATH')
    if gcs_path is None:
       gcs_path = daily_params['GCS_DAILY_PATH'].format(
                      version=version)

    mfest_commit = conf.get('MFEST_COMMIT')
    if mfest_commit is None:
       timestamp = time.mktime(date.timetuple())
       mfest_commit = daily_params['MFEST_COMMIT'].format(
          branch=branch, timestamp=timestamp)

    default_conf = environment_config.get_default_config(
        branch=branch,
        gcs_path=gcs_path,
        mfest_commit=mfest_commit,
        verify_consistency=daily_params['VERIFY_CONSISTENCY'],
        version=version)

    config_settings = dict()
    for name in default_conf.iterkeys():
      config_settings[name] = conf.get(name) or default_conf[name]

    return config_settings

  dag_name = 'istio_daily_' + branch
  dag, tasks = MakeCommonDag(
       DailyGenerateTestArgs,
       name=dag_name, schedule_interval='15 9 * * *')
  tasks['mark_daily_complete'] = MakeMarkComplete(dag)

  #tasks['generate_workflow_args']
  tasks['get_git_commit'                 ].set_upstream(tasks['generate_workflow_args'])
  tasks['run_cloud_builder'              ].set_upstream(tasks['get_git_commit'])
  tasks['run_release_qualification_tests'].set_upstream(tasks['run_cloud_builder'])
  tasks['copy_files_for_release'         ].set_upstream(tasks['run_release_qualification_tests'])
  tasks['mark_daily_complete'            ].set_upstream(tasks['copy_files_for_release'])

  return dag



def MonthlyPipeline(branch):
  MONTHLY_RELEASE_TRIGGER = '15 17 * * 4#3'

  def MonthlyGenerateTestArgs(**kwargs):
    release_params = dict(
      DOCKER_HUB='istio',
      GCR_RELEASE_DEST='istio-io',
      GCS_GITHUB_PATH='istio-secrets/github.txt.enc',
      GCS_MONTHLY_RELEASE_PATH='istio-release/releases/{version}',
      RELEASE_PROJECT_ID='istio-io',
      VERIFY_CONSISTENCY='true')

    """Loads the configuration that will be used for this Iteration."""
    conf = kwargs['dag_run'].conf
    if conf is None:
      conf = dict()

    # If version is overriden then we should use it otherwise we use it's
    # default or monthly value.
    version = conf.get('VERSION') or AirflowGetVariableOrBaseCase(branch+'-version', None)
    if not version or version == 'INVALID':
      raise ValueError('version needs to be provided')
    Variable.set(branch+'-version', 'INVALID')

    GCS_MONTHLY_STAGE_PATH='prerelease/{version}'
    gcs_path = GCS_MONTHLY_STAGE_PATH.format(version=version)

    mfest_commit = conf.get('MFEST_COMMIT')
    if mfest_commit is None:
      mfest_commit = branch

    default_conf = environment_config.get_default_config(
        branch=branch,
        gcs_path=gcs_path,
        mfest_commit=mfest_commit,
        verify_consistency=release_params['VERIFY_CONSISTENCY'],
        version=version)

    config_settings = dict()
    for name in default_conf.iterkeys():
      config_settings[name] = conf.get(name) or default_conf[name]

    monthly_conf = dict(release_params)
    monthly_conf['GCS_MONTHLY_RELEASE_PATH'] = release_params['GCS_MONTHLY_RELEASE_PATH'].format(
		version=config_settings['VERSION'])
    for name in release_params.iterkeys():
      config_settings[name] = conf.get(name) or monthly_conf[name]

    return config_settings

  dag, tasks = MakeCommonDag(
    MonthlyGenerateTestArgs,
    'istio_monthly_'+branch,
    schedule_interval=MONTHLY_RELEASE_TRIGGER)

  release_push_github_docker_template = """
{% set m_commit = task_instance.xcom_pull(task_ids='get_git_commit') %}
{% set settings = task_instance.xcom_pull(task_ids='generate_workflow_args') %}
gsutil cp gs://{{ settings.GCS_RELEASE_TOOLS_PATH }}/data/release/*.json .
gsutil cp gs://{{ settings.GCS_RELEASE_TOOLS_PATH }}/data/release/*.sh .
chmod u+x *
./start_gcb_publish.sh \
-p "{{ settings.RELEASE_PROJECT_ID }}" -a "{{ settings.SVC_ACCT }}"  \
-v "{{ settings.VERSION }}" -s "{{ settings.GCS_FULL_STAGING_PATH }}" \
-b "{{ settings.GCS_MONTHLY_RELEASE_PATH }}" -r "{{ settings.GCR_RELEASE_DEST }}" \
-g "{{ settings.GCS_GITHUB_PATH }}" -u "{{ settings.MFEST_URL }}" \
-t "{{ m_commit }}" -m "{{ settings.MFEST_FILE }}" \
-h "{{ settings.GITHUB_ORG }}" -i "{{ settings.GITHUB_REPO }}" \
-d "{{ settings.DOCKER_HUB}}" -w
"""

  github_and_docker_release = BashOperator(
    task_id='github_and_docker_release',
    bash_command=release_push_github_docker_template,
    dag=dag)
  tasks['github_and_docker_release'] = github_and_docker_release

  release_tag_github_template = """
{% set m_commit = task_instance.xcom_pull(task_ids='get_git_commit') %}
{% set settings = task_instance.xcom_pull(task_ids='generate_workflow_args') %}
gsutil cp gs://{{ settings.GCS_RELEASE_TOOLS_PATH }}/data/release/*.json .
gsutil cp gs://{{ settings.GCS_RELEASE_TOOLS_PATH }}/data/release/*.sh .
chmod u+x *
./start_gcb_tag.sh \
-p "{{ settings.RELEASE_PROJECT_ID }}" \
-h "{{ settings.GITHUB_ORG }}" -a "{{ settings.SVC_ACCT }}"  \
-v "{{ settings.VERSION }}"   -e "istio_releaser_bot@example.com" \
-n "IstioReleaserBot" -s "{{ settings.GCS_FULL_STAGING_PATH }}" \
-g "{{ settings.GCS_GITHUB_PATH }}" -u "{{ settings.MFEST_URL }}" \
-t "{{ m_commit }}" -m "{{ settings.MFEST_FILE }}" -w
"""

  github_tag_repos = BashOperator(
    task_id='github_tag_repos',
    bash_command=release_tag_github_template,
    dag=dag)
  tasks['github_tag_repos'] = github_tag_repos


  def ReportMonthlySuccessful(task_instance, **kwargs):
    del kwargs

  mark_monthly_complete = PythonOperator(
    task_id='mark_monthly_complete',
    python_callable=ReportMonthlySuccessful,
    provide_context=True,
    dag=dag,
  )
  tasks['mark_monthly_complete'] = mark_monthly_complete

# tasks['generate_workflow_args']
  tasks['get_git_commit'                 ].set_upstream(tasks['generate_workflow_args'])
  tasks['run_cloud_builder'              ].set_upstream(tasks['get_git_commit'])
  tasks['run_release_qualification_tests'].set_upstream(tasks['run_cloud_builder'])
  tasks['copy_files_for_release'         ].set_upstream(tasks['run_release_qualification_tests'])
  tasks['github_and_docker_release'      ].set_upstream(tasks['copy_files_for_release'])
  tasks['github_tag_repos'               ].set_upstream(tasks['github_and_docker_release'])
  tasks['mark_monthly_complete'          ].set_upstream(tasks['github_tag_repos'])

  return dag
