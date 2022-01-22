package containers

import (
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/rs/zerolog/log"
)

// Inspect returns inspect resuls of the specific container
func Inspect(id string) (string, error) {
	log.Debug().Msgf("pdcs: podman container inspect %s", id)
	var report string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := containers.Inspect(conn, id, new(containers.InspectOptions))
	if err != nil {
		return report, err
	}

	report, err = utils.GetJSONOutput(response)
	if err != nil {
		return report, err
	}
	log.Debug().Msgf("pdcs: %s", report)
	return report, nil
}
