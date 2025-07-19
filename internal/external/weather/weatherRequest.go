package weather

type (
	Requestor interface {
		build()
		send()
	}

	Request struct {
		Requestor
	}
)
