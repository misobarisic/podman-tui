#!/usr/bin/env bats
#
# podman-tui pods view functionality tests
#

load helpers
load helpers_tui

@test "pod create" {
    podman pod rm $TEST_POD_NAME || echo done
    podman image pull pause:3.5 || echo done
    podman network create podman || echo done

    # switch to pods view
    # select create command from pod commands dialog
    # fillout name field
    # fillout label field
    # switch to "networking" create view
    # go to networks dropdown widget and select network name from available networks
    # go to "Create" button and press Enter
    podman_tui_set_view "pods"
    podman_tui_select_pod_cmd "create"
    podman_tui_send_inputs "Enter"
    podman_tui_send_inputs $TEST_POD_NAME "Tab" "Tab" $TEST_LABEL
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Down" "Down" "Down"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Tab" "Down" "Down" "Enter"
    podman_tui_send_inputs "Tab" "Tab" 
    podman_tui_send_inputs "Enter"
    sleep 4

    run_podman pod ls --format "'{{ .Name }}'"
    assert "$output" "=~" "$TEST_POD_NAME" "expected $TEST_POD_NAME to be in the list"

    run_podman pod inspect $TEST_POD_NAME --format "'{{ json .InfraConfig.Networks }}'"
    assert "$output" =~ "podman" "expected $TEST_POD_NAME to have podman network"

}

@test "pod start" {
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 0 | grep "$TEST_POD_NAME" | awk '{print $1}')

    # switch to pods view
    # select test pod from list
    # select start command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item $pod_index
    podman_tui_select_pod_cmd "start"
    sleep 2

    run_podman pod ls --filter="name=$TEST_POD_NAME" --format "'{{ .Status}}'"
    assert $output =~ "Running" "expected $TEST_POD_NAME running"
}

@test "pod pause" {
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 0 | grep "$TEST_POD_NAME" | awk '{print $1}')

    # switch to pods view
    # select test pod from list
    # select pause command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item $pod_index
    podman_tui_select_pod_cmd "pause"
    sleep 2

    run_podman pod ls --filter="name=$TEST_POD_NAME" --format "'{{ .Status}}'"
    assert $output =~ "Paused" "expected $TEST_POD_NAME running"
}

@test "pod unpause" {
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 0 | grep "$TEST_POD_NAME" | awk '{print $1}')

    # switch to pods view
    # select test pod from list
    # select unpause command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item $pod_index
    podman_tui_select_pod_cmd "unpause"
    sleep 2

    run_podman pod ls --filter="name=$TEST_POD_NAME" --format "'{{ .Status}}'"
    assert $output =~ "Running" "expected $TEST_POD_NAME running"
}

@test "pod stop" {
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 0 | grep "$TEST_POD_NAME" | awk '{print $1}')

    # switch to pods view
    # select test pod from list
    # select stop command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item $pod_index
    podman_tui_select_pod_cmd "stop"
    sleep 2

    run_podman pod ls --filter="name=$TEST_POD_NAME" --format "'{{ .Status}}'"
    assert $output =~ "Exited" "expected $TEST_POD_NAME exited"
}

@test "pod restart" {
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 0 | grep "$TEST_POD_NAME" | awk '{print $1}')

    # switch to pods view
    # select test pod from list
    # select restart command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item $pod_index
    podman_tui_select_pod_cmd "restart"
    sleep 2

    run_podman pod ls --filter="name=$TEST_POD_NAME" --format "'{{ .Status}}'"
    assert $output =~ "Running" "expected $TEST_POD_NAME exited"
}

@test "pod kill" {
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 0 | grep "$TEST_POD_NAME" | awk '{print $1}')

    # switch to pods view
    # select test pod from list
    # select kill command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item $pod_index
    podman_tui_select_pod_cmd "kill"
    sleep 2

    run_podman pod ls --filter="name=$TEST_POD_NAME" --format "'{{ .Status}}'"
    assert $output =~ "Exited" "expected $TEST_POD_NAME exited"
}

@test "pod inspect" {
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 0 | grep "$TEST_POD_NAME" | awk '{print $1}')
    
    # switch to pods view
    # select test pod from list
    # select inspect command from pod commands dialog
    # close pod inspect result message dialog
    podman_tui_set_view "pods"
    podman_tui_select_item $pod_index
    podman_tui_select_pod_cmd "inspect"
    sleep 2
    podman_tui_send_inputs "Enter"

    run_helper sed -n '/  "Labels": {/, /  },/p' $PODMAN_TUI_LOG
    assert "$output" "=~" "\"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\"" "expected \"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\" in pod inspect"
}

@test "pod remove" {
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 0 | grep "$TEST_POD_NAME" | awk '{print $1}')

    # switch to pods view
    # select test pod from list
    # select remove command from pod commands dialog
    # confirm the operation on warnings dialog
    podman_tui_set_view "pods"
    podman_tui_select_item $pod_index
    podman_tui_select_pod_cmd "remove"
    podman_tui_send_inputs "Enter"
    podman_tui_send_inputs "Enter"
    sleep 2

    run_podman pod ls --format "{{ .Name }}" --filter "name=$TEST_POD_NAME"
    assert "$output" "=~" "" "expected $TEST_POD_NAME pod removal"
}

@test "pod prune" {
    podman pod create --name $TEST_POD_NAME --label $TEST_LABEL || echo done
    sleep 2
    
    # switch to pods view
    # select prune command from pod commands dialog
    # confirm the operation on warnings dialog
    podman_tui_set_view "pods"
    podman_tui_select_pod_cmd "prune"
    podman_tui_send_inputs "Enter"
    sleep 3

    run_podman pod ls --format "{{ .Name }}" --filter "name=$TEST_POD_NAME"
    assert "$output" "=~" "" "expected at least $TEST_POD_NAME pod removal"
}
