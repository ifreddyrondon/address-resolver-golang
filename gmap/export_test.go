package gmap

func SetClientInstance(client AddressClient) {
	serviceInstance = &AddressService{client: client}
}
