package libresponse

// ErrorCommitGeneral -
func (r *Response) ErrorCommitGeneral() *Response {
	r.Code = 500
	r.Message = "common.error.server.internal"
	r.Error = "common.error.service.commit.general"
	return r
}

// ErrorTxBegin -
func (r *Response) ErrorTxBegin() *Response {
	r.Code = 500
	r.Message = "common.error.server.internal"
	r.Error = "common.error.service.tx.begin"
	return r
}

// ErrorTxCreate -
func (r *Response) ErrorTxCreate() *Response {
	r.Code = 500
	r.Message = "common.error.server.internal"
	r.Error = "common.error.service.tx.create"
	return r
}

// ErrorCommitCreate -
func (r *Response) ErrorCommitCreate() *Response {
	r.Code = 500
	r.Message = "common.error.server.internal"
	r.Error = "common.error.service.commit.create"
	return r
}

// ErrorTxUpdate -
func (r *Response) ErrorTxUpdate() *Response {
	r.Code = 500
	r.Message = "common.error.server.internal"
	r.Error = "common.error.service.tx.update"
	return r
}

// ErrorCommitUpdate -
func (r *Response) ErrorCommitUpdate() *Response {
	r.Code = 500
	r.Message = "common.error.server.internal"
	r.Error = "common.error.service.commit.update"
	return r
}

// ErrorTxDelete -
func (r *Response) ErrorTxDelete() *Response {
	r.Code = 500
	r.Message = "common.error.server.internal"
	r.Error = "common.error.service.tx.delete"
	return r
}

// ErrorCommitDelete -
func (r *Response) ErrorCommitDelete() *Response {
	r.Code = 500
	r.Message = "common.error.server.internal"
	r.Error = "common.error.service.commit.delete"
	return r
}
