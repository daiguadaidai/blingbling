package reviewer

type Reviewer interface {
	Review() (*ReviewMSG)
}

