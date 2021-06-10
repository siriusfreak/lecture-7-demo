package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/afero"
	"github.com/yandex/pandora/cli"
	"github.com/yandex/pandora/components/phttp/import"
	"github.com/yandex/pandora/core"
	"github.com/yandex/pandora/core/aggregator/netsample"
	"github.com/yandex/pandora/core/import"
	"github.com/yandex/pandora/core/register"

	"log"

	"google.golang.org/grpc"

	pb "gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo"
)

type Ammo struct {
	Tag         string
	Param1      string
	Param2      string
	Param3      string
}

type Sample struct {
	URL              string
	ShootTimeSeconds float64
}

type GunConfig struct {
	Target string `validate:"required"` // Configuration will fail, without target defined
}

type Gun struct {
	// Configured on construction.
	client grpc.ClientConn
	conf   GunConfig
	// Configured on Bind, before shooting
	aggr core.Aggregator // May be your custom Aggregator.
	core.GunDeps
}

func NewGun(conf GunConfig) *Gun {
	return &Gun{conf: conf}
}

func (g *Gun) Bind(aggr core.Aggregator, deps core.GunDeps) error {
	conn, err := grpc.Dial(
		g.conf.Target,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithUserAgent("pandora load test"))
	if err != nil {
		log.Fatalf("FATAL: %s", err)
	}
	g.client = *conn
	g.aggr = aggr
	g.GunDeps = deps
	return nil
}

func (g *Gun) Shoot(ammo core.Ammo) {
	customAmmo := ammo.(*Ammo) // Shoot will panic on unexpected ammo type. Panic cancels shooting.
	g.shoot(customAmmo)
}

func (g *Gun) shoot(ammo *Ammo) {
	//start := time.Now()
	code := 0
	sample := netsample.Acquire(ammo.Tag)

	conn := g.client
	client := pb.NewLecture7DemoClient(&conn)

	// prepare list of ids from ammo
	var itemIDs []int64
	for _, id := range strings.Split(ammo.Param1, ",") {
		if id == "" {
			continue
		}
		itemID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		itemIDs = append(itemIDs, itemID)
	}

	// make scores request
	req := http.Request{}
	out, err := client.AddV1(req.Context(), &pb.AddRequestV1{
		Id:          1,
		Text:        "test",
		Result:      true,
		CallbackUrl: "",
	})

	if err != nil {
		code = 0
	}

	if out != nil {
		code = 200
	}

	defer func() {
		sample.SetProtoCode(code)
		g.aggr.Report(sample)
	}()
}

func main() {
	//debug.SetGCPercent(-1)
	// Standard imports.
	fs := afero.NewOsFs()
	coreimport.Import(fs)
	// May not be imported, if you don't need http guns and etc.
	phttp.Import(fs)

	// Custom imports. Integrate your custom types into configuration system.
	coreimport.RegisterCustomJSONProvider("my-custom-provider-name", func() core.Ammo { return &Ammo{} })

	register.Gun("my-custom-gun-name", NewGun, func() GunConfig {
		return GunConfig{
			Target: "default target",
		}
	})


	cli.Run()
}