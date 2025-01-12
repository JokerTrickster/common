package schema

type Tokens struct {
	ID               int    `json:"id"`
	UserID           int    `json:"userID"`
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	RefreshExpiredAt int64  `json:"refreshExpiredAt"`
}

type Users struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Score    int    `json:"score"`
	State    string `json:"state"`
	RoomID   int    `json:"roomID"`
	Provider string `json:"provider"` // google, kakao, naver, email
}
type Rooms struct {
	ID           int    `json:"id"`
	CurrentCount int    `json:"currentCount"` //방 현재 인원
	MaxCount     int    `json:"maxCount"`     //방 최대 인원
	MinCount     int    `json:"minCount"`     //방 최소 인원
	Name         string `json:"name"`         //방 이름
	Password     string `json:"password"`     //방 비밀번호 (옵셔널))
	State        string `json:"state"`        //방 상태 (대기, 진행, 종료)
	OwnerID      int    `json:"ownerID"`      //방 주인 아이디
}
