package core

type JWT struct {
	Token        string `json:"jwt_token" binding:"required"`
	IsRegistered bool   `json:"-"`
	JwtRole      string `json:"-"`
}

func (t *JWT) ToResponse() JWTResponse {
	return JWTResponse{Token: t.Token, IsRegistered: t.IsRegistered, Role: t.JwtRole}
}

type JWTResponse struct {
	Token        string `json:"jwt_token" binding:"required"`
	IsRegistered bool   `json:"is_registered"`
	Role         string `json:"role"`
}
