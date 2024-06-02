package nomad

import "github.com/hashicorp/nomad/api"

func New() *api.Client {
	client, err := api.NewClient(api.DefaultConfig())

	if err != nil {
		panic(err)
	}

	return client
}
