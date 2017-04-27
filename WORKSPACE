workspace(name = "com_github_istio_istio")

git_repository(
    name = "io_bazel_rules_go",
    commit = "87cdda3fc0fd65c63ef0316533be03ea4956f809",  # April 7 2017 (0.4.2)
    remote = "https://github.com/bazelbuild/rules_go.git",
)

load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "new_go_repository")

go_repositories()

git_repository(
    name = "org_pubref_rules_protobuf",
    commit = "9ede1dbc38f0b89ae6cd8e206a22dd93cc1d5637",  # Mar 31 2017 (gogo* support)
    remote = "https://github.com/pubref/rules_protobuf",
)

load("@org_pubref_rules_protobuf//gogo:rules.bzl", "gogo_proto_repositories")
load("@org_pubref_rules_protobuf//cpp:rules.bzl", "cpp_proto_repositories")

cpp_proto_repositories()

gogo_proto_repositories()

new_go_repository(
    name = "com_github_golang_glog",
    commit = "23def4e6c14b4da8ac2ed8007337bc5eb5007998",  # Jan 26, 2016 (no releases)
    importpath = "github.com/golang/glog",
)

new_go_repository(
    name = "com_google_cloud_go",
    commit = "57377bad3486b37af17b47230a61603794c798ae",
    importpath = "cloud.google.com/go",
)

new_go_repository(
    name = "org_golang_x_net",
    commit = "242b6b35177ec3909636b6cf6a47e8c2c6324b5d",
    importpath = "golang.org/x/net",
)

new_go_repository(
    name = "org_golang_x_oauth2",
    commit = "314dd2c0bf3ebd592ec0d20847d27e79d0dbe8dd",
    importpath = "golang.org/x/oauth2",
)

new_go_repository(
    name = "org_golang_google_api",
    commit = "48e49d1645e228d1c50c3d54fb476b2224477303",
    importpath = "google.golang.org/api",
)

new_go_repository(
    name = "org_golang_google_grpc",
    commit = "377586b314e142ce186a0644138c92fe55b9162e",
    importpath = "google.golang.org/grpc",
)

new_go_repository(
    name = "org_golang_google_genproto",
    commit = "411e09b969b1170a9f0c467558eb4c4c110d9c77",
    importpath = "google.golang.org/genproto",
)

new_go_repository(
    name = "com_github_googleapis_gax_go",
    commit = "9af46dd5a1713e8b5cd71106287eba3cefdde50b",
    importpath = "github.com/googleapis/gax-go",
)

new_go_repository(
    name = "com_github_google_uuid",
    commit = "6a5e28554805e78ea6141142aba763936c4761c0",
    importpath = "github.com/google/uuid",
)

new_go_repository(
    name = "com_github_golang_protobuf",
    commit = "2bba0603135d7d7f5cb73b2125beeda19c09f4ef",
    importpath = "github.com/golang/protobuf",
)

new_go_repository(
    name = "com_github_pmezard_go_difflib",
    commit = "d8ed2627bdf02c080bf22230dbb337003b7aba2d",
    importpath = "github.com/pmezard/go-difflib",
)

new_go_repository(
    name = "com_github_hashicorp_errwrap",
    commit = "7554cd9344cec97297fa6649b055a8c98c2a1e55",
    importpath = "github.com/hashicorp/errwrap",
)

new_go_repository(
    name = "com_github_hashicorp_go_multierror",
    commit = "8484912a3b9987857bac52e0c5fec2b95f419628",
    importpath = "github.com/hashicorp/go-multierror",
)
