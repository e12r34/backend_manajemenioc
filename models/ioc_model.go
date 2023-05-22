package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ioc struct {
	Id          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Type        string             `json:"type,omitempty" validate:"required"`
	Indicator   string             `json:"ioc,omitempty" validate:"required"`
	Description string             `json:"description,omitempty"`
	Source      string             `json:"source,omitempty"`
	OTX_Hash    string             `json:"otx_hash,omitempty"`
}

type Ioc_Bson struct {
	Id          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Type        string             `bson:"type,omitempty" validate:"required"`
	Indicator   string             `bson:"ioc,omitempty" validate:"required"`
	Description string             `bson:"description,omitempty"`
	Source      string             `bson:"source,omitempty"`
	OTX_Hash    string             `bson:"otx_hash,omitempty"`
}

type Ioc_input struct {
	Type        string `json:"type,omitempty" validate:"required"`
	Indicator   string `json:"ioc,omitempty" validate:"required"`
	Description string `json:"description,omitempty"`
	Source      string `json:"source,omitempty"`
	OTX_Hash    string `json:"otx_hash,omitempty"`
}

type Ioc_Edit struct {
	Type        string `json:"type,omitempty"`
	Indicator   string `bson:"indicator" json:"ioc,omitempty"`
	Description string `json:"description,omitempty"`
	Source      string `json:"source,omitempty"`
	OTX_Hash    string `json:"otx_hash,omitempty"`
}

type Ioc_Get struct {
	Id []primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

type Ioc_Thor struct {
	Indicator string `bson:"indicator,omitempty" json:"ioc,omitempty"`
	Id        string `bson:"otx_id,omitempty" json:"id,omitempty"`
	Content   string `bson:"content,omitempty"`
	Name      string
}

type Ioc_MG_BAE struct {
	Type      string `bson:"type" json:"type"`
	Indicator string `bson:"indicator" json:"-"`
	Id        string `bson:"otx_id" json:"-"`
}
