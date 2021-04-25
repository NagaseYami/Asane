package system

type IRecivedMessageObject interface {
	Reply(IIdentityDocument, IIdentityDocument)
	GetSenderID() IIdentityDocument
	GetGroupID() IIdentityDocument
	GetRawMessage() string
}

type IIdentityDocument interface {
	GetStringID() string
	GetInt64ID() int64
	GetInt32ID() int32
}
