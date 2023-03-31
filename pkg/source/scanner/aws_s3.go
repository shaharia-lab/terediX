package scanner

import "teredix/pkg/resource"

type AWSS3 struct {
	SourceName   string
	AccessKey    string
	SecretKey    string
	SessionToken string
	Zone         string
}

func NewAWSS3(sourceName string, accessKey string, secretKey string, sessionToken string, zone string) *AWSS3 {
	return &AWSS3{
		SourceName:   sourceName,
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		SessionToken: sessionToken,
		Zone:         zone,
	}
}

func (a AWSS3) Scan() []resource.Resource {
	return []resource.Resource{}
}
