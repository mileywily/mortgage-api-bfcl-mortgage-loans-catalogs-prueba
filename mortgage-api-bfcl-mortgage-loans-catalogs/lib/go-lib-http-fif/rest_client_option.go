package http_fif

type RestClientOption func(restClient) restClient

func Encoder(encoder DataEncoder) RestClientOption {
	return func(client restClient) restClient {
		client.encoder = encoder
		return client
	}
}
