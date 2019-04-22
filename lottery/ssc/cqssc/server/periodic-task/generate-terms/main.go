package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"

	// "encoding/json"

	"github.com/golang/protobuf/ptypes"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "audi/lottery/ssc/cqssc/proto"
	_ "audi/pkg/logger"
)

var (
	addr string
	wg   sync.WaitGroup
)

func init() {
	flag.StringVar(&addr, "addr", "lottery.ssc.cqssc:3721", "cqssc grpc server address")

	flag.Parse()
}

type Period struct {
	Id    int64     `json:"id"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type Durantion struct {
	Start    string
	Duration time.Duration
}

var Durations = []Durantion{
	Durantion{"00:10", 20 * time.Minute},
	Durantion{"00:30", 20 * time.Minute},
	Durantion{"00:50", 20 * time.Minute},
	Durantion{"01:10", 20 * time.Minute},
	Durantion{"01:30", 20 * time.Minute},
	Durantion{"01:50", 20 * time.Minute},
	Durantion{"02:10", 20 * time.Minute},
	Durantion{"02:30", 20 * time.Minute},
	Durantion{"02:50", 20 * time.Minute},

	Durantion{"07:10", 20 * time.Minute},
	Durantion{"07:30", 20 * time.Minute},
	Durantion{"07:50", 20 * time.Minute},
	Durantion{"08:10", 20 * time.Minute},
	Durantion{"08:30", 20 * time.Minute},
	Durantion{"08:50", 20 * time.Minute},
	Durantion{"09:10", 20 * time.Minute},
	Durantion{"09:30", 20 * time.Minute},
	Durantion{"09:50", 20 * time.Minute},
	Durantion{"10:10", 20 * time.Minute},
	Durantion{"10:30", 20 * time.Minute},
	Durantion{"10:50", 20 * time.Minute},
	Durantion{"11:10", 20 * time.Minute},
	Durantion{"11:30", 20 * time.Minute},
	Durantion{"11:50", 20 * time.Minute},
	Durantion{"12:10", 20 * time.Minute},
	Durantion{"12:30", 20 * time.Minute},
	Durantion{"12:50", 20 * time.Minute},
	Durantion{"13:10", 20 * time.Minute},
	Durantion{"13:30", 20 * time.Minute},
	Durantion{"13:50", 20 * time.Minute},
	Durantion{"14:10", 20 * time.Minute},
	Durantion{"14:30", 20 * time.Minute},
	Durantion{"14:50", 20 * time.Minute},
	Durantion{"15:10", 20 * time.Minute},
	Durantion{"15:30", 20 * time.Minute},
	Durantion{"15:50", 20 * time.Minute},
	Durantion{"16:10", 20 * time.Minute},
	Durantion{"16:30", 20 * time.Minute},
	Durantion{"16:50", 20 * time.Minute},
	Durantion{"17:10", 20 * time.Minute},
	Durantion{"17:30", 20 * time.Minute},
	Durantion{"17:50", 20 * time.Minute},
	Durantion{"18:10", 20 * time.Minute},
	Durantion{"18:30", 20 * time.Minute},
	Durantion{"18:50", 20 * time.Minute},
	Durantion{"19:10", 20 * time.Minute},
	Durantion{"19:30", 20 * time.Minute},
	Durantion{"19:50", 20 * time.Minute},
	Durantion{"20:10", 20 * time.Minute},
	Durantion{"20:30", 20 * time.Minute},
	Durantion{"20:50", 20 * time.Minute},
	Durantion{"21:10", 20 * time.Minute},
	Durantion{"21:30", 20 * time.Minute},
	Durantion{"21:50", 20 * time.Minute},
	Durantion{"22:10", 20 * time.Minute},
	Durantion{"22:30", 20 * time.Minute},
	Durantion{"22:50", 20 * time.Minute},
	Durantion{"23:10", 20 * time.Minute},
	Durantion{"23:30", 20 * time.Minute},
}

func convertToPeriod(dateStr string, idx int, d Durantion) *Period {
	id, err := strconv.ParseInt(fmt.Sprintf("%s%0#3d", dateStr, idx+1), 10, 64)
	if err != nil {
		panic(err)
	}

	start, err := time.Parse("20060102 15:04", fmt.Sprintf("%s %s", dateStr, d.Start))
	if err != nil {
		panic(err)
	}
	end := start.Add(d.Duration)
	period := &Period{Id: id, Start: start, End: end}

	return period
}

func GetPeriodsByDate(dateStr string) []Period {
	var ps []Period
	for seq, d := range Durations {
		p := convertToPeriod(dateStr, seq, d)
		ps = append(ps, *p)

	}
	return ps
}

func createTermFunc(ctx context.Context, c pb.CqsscServiceClient, req *pb.CreateTermReq) {
	logger := log.WithField("id", req.Id)
	if _, err := c.CreateTerm(ctx, req); err != nil {
		logger.Error(err)
	} else {
		logger.Info("create term success")
	}

	wg.Done()
	return
}

func createTermJob() {
	log.Info("Create terms start")

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Error(err)
	}
	defer conn.Close()

	client := pb.NewCqsscServiceClient(conn)
	ctx := context.Background()

	// tomorrow
	tomrStr := time.Now().AddDate(0, 0, 1).Format("20060102")
	periods := GetPeriodsByDate(tomrStr)
	for _, p := range periods {
		wg.Add(1)
		// fmt.Printf("%v \n", p)
		st, _ := ptypes.TimestampProto(p.Start)
		et, _ := ptypes.TimestampProto(p.End)
		req := pb.CreateTermReq{Id: p.Id, StartFrom: st, EndTo: et}
		go createTermFunc(ctx, client, &req)
	}
	wg.Wait()
	
	log.Info("Create terms completed")
}

func main() {
	p := "0 0 4 * * *"
	c := cron.New()
	c.AddFunc(p, func() {
		go createTermJob()
	})
	c.Start()

	log.Info("Task starting")
	log.Infof("Create ters on: %s", p)

	select {}
}
