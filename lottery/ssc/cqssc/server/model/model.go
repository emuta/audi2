package model

import (
	"time"

	"github.com/lib/pq"
	"github.com/jinzhu/gorm"
)

type Config struct {
    Id     int64   `gorm:"primary_key:true"`
    Name   string
    Tag    string
    Odds   float64
    Comm   float64
    Price  float64
    Active bool

}

func (Config) TableName() string {
    return "cqssc.config"
}


type Unit struct {
    Id    int64   `gorm:"primary_key:true"`
    Name  string
    Value float64
}

func (Unit) TableName() string {
    return "cqssc.unit"
}


type Catg struct {
    Id   int32  `gorm:"primary_key:true"`
    Name string
    Tag  string
    Pref bool
}

func (Catg) TableName() string {
    return "cqssc.catg"
}

type Group struct {
    Id   int32  `gorm:"primary_key:true"`
    Name string
    Tag  string
}

func (Group) TableName() string {
    return "cqssc.group"
}

type Play struct {
    Id        int32  `gorm:"primary_key:true"`
    Name      string
    Tag       string
    Pref      bool
    Active    bool
    Pr        int32
    CatgId    int32
    GroupId   int32 `gorm:"column:group_id"`
    // Units     []int32
    Units     pq.Int64Array
}

func (Play) TableName() string {
    return "cqssc.play"
}

type Term struct {
    Id        int64 `gorm:"primary_key:true"`
    StartFrom time.Time
    EndTo     time.Time
    Codes     pq.StringArray
    OpenedAt  time.Time
    SettledAt pq.NullTime
    RevokedAt pq.NullTime
}

func (Term) TableName() string {
    return "cqssc.term"
}

type Bet struct {
	Id        int64          `gorm:"primary_key:true"`
	CreatedAt time.Time
	UserId    int64
	Odds      float64
	PlayId    int32
	UnitId    int64
	Comm      float64
	ChaseStop bool
	Codes     pq.StringArray

	// Plans     []Plan `gorm:"-"`
}

func (Bet) TableName() string {
	return "cqssc.bet"
}

type BetPlan struct {
	Id      int64     `gorm:"primary_key:true"`
	BetId   int64
	TermId  int64
	Times   int64
	Qty     int64
	Payment float64
	Rebate  float64
	Bonus   float64
	// Stats   PlanStats `gorm:"-"`
}

func (BetPlan) TableName() string {
	return "cqssc.bet_plan"
}

func (p *BetPlan) AfterCreate(tx *gorm.DB) error {
	stats := BetPlanStats{
		Id:      p.Id,
		BetId:   p.BetId,
		Settled: false,
		Revoked: false,
		Payment: p.Payment,
		Bonus:   0,
		Rebate:  0,
		Win:     0,
	}
	return tx.Create(&stats).Error
}

type BetPlanStats struct {
	Id        int64       `gorm:"primary_key:true"`
	BetId     int64
	Settled   bool
	SettledAt pq.NullTime
	Revoked   bool
	RevokedAt pq.NullTime

	Win       int64
	Payment   float64
	Bonus     float64
	Rebate    float64
}

func (BetPlanStats) TableName() string {
	return "cqssc.bet_plan_stats"
}

func (s *BetPlanStats) AfterUpdate(tx *gorm.DB) error {
	// should refresh bet stats
	return nil
}