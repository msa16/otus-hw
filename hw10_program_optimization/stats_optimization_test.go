//go:build bench
// +build bench

package hw10programoptimization

import (
	"archive/zip"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	mb          uint64 = 1 << 20
	memoryLimit uint64 = 30 * mb

	timeLimit = 300 * time.Millisecond
)

// go test -v -count=1 -timeout=30s -tags bench .
// бенчмарки по скорости и памяти
// go test -bench=. -benchmem .
// запись профилировки
// go test -bench=. -benchmem -cpuprofile=cpu.out -memprofile=mem.out .
// просмотр профилировки
// go tool pprof -http=":8090" 01_bench.test mem.out
func TestGetDomainStat_Time_And_Memory(t *testing.T) {
	bench := func(b *testing.B) {
		b.Helper()
		b.StopTimer()

		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(t, err)
		defer r.Close()

		require.Equal(t, 1, len(r.File))

		data, err := r.File[0].Open()
		require.NoError(t, err)

		b.StartTimer()
		stat, err := GetDomainStat(data, "biz")
		b.StopTimer()
		require.NoError(t, err)

		require.Equal(t, expectedBizStat, stat)
	}

	result := testing.Benchmark(bench)
	mem := result.MemBytes
	t.Logf("time used: %s / %s", result.T, timeLimit)
	t.Logf("memory used: %dMb / %dMb", mem/mb, memoryLimit/mb)

	require.Less(t, int64(result.T), int64(timeLimit), "the program is too slow")
	require.Less(t, mem, memoryLimit, "the program is too greedy")
}

var expectedBizStat = DomainStat{
	"abata.biz":         25,
	"abatz.biz":         25,
	"agimba.biz":        28,
	"agivu.biz":         17,
	"aibox.biz":         31,
	"ailane.biz":        23,
	"aimbo.biz":         25,
	"aimbu.biz":         36,
	"ainyx.biz":         35,
	"aivee.biz":         25,
	"avamba.biz":        21,
	"avamm.biz":         17,
	"avavee.biz":        35,
	"avaveo.biz":        30,
	"babbleblab.biz":    29,
	"babbleopia.biz":    36,
	"babbleset.biz":     28,
	"babblestorm.biz":   29,
	"blognation.biz":    32,
	"blogpad.biz":       34,
	"blogspan.biz":      21,
	"blogtag.biz":       23,
	"blogtags.biz":      34,
	"blogxs.biz":        35,
	"bluejam.biz":       36,
	"bluezoom.biz":      27,
	"brainbox.biz":      30,
	"brainlounge.biz":   38,
	"brainsphere.biz":   31,
	"brainverse.biz":    39,
	"brightbean.biz":    23,
	"brightdog.biz":     32,
	"browseblab.biz":    31,
	"browsebug.biz":     25,
	"browsecat.biz":     34,
	"browsedrive.biz":   24,
	"browsetype.biz":    34,
	"browsezoom.biz":    29,
	"bubblebox.biz":     19,
	"bubblemix.biz":     38,
	"bubbletube.biz":    34,
	"buzzbean.biz":      26,
	"buzzdog.biz":       30,
	"buzzshare.biz":     26,
	"buzzster.biz":      28,
	"camido.biz":        27,
	"camimbo.biz":       36,
	"centidel.biz":      32,
	"centimia.biz":      17,
	"centizu.biz":       18,
	"chatterbridge.biz": 30,
	"chatterpoint.biz":  32,
	"cogibox.biz":       30,
	"cogidoo.biz":       34,
	"cogilith.biz":      24,
	"dabfeed.biz":       26,
	"dabjam.biz":        30,
	"dablist.biz":       30,
	"dabshots.biz":      33,
	"dabtype.biz":       21,
	"dabvine.biz":       26,
	"dabz.biz":          19,
	"dazzlesphere.biz":  24,
	"demimbu.biz":       27,
	"demivee.biz":       39,
	"demizz.biz":        30,
	"devbug.biz":        20,
	"devcast.biz":       35,
	"devify.biz":        27,
	"devpoint.biz":      26,
	"devpulse.biz":      27,
	"devshare.biz":      30,
	"digitube.biz":      30,
	"divanoodle.biz":    33,
	"divape.biz":        32,
	"divavu.biz":        28,
	"dynabox.biz":       66,
	"dynava.biz":        21,
	"dynazzy.biz":       29,
	"eabox.biz":         28,
	"eadel.biz":         25,
	"eamia.biz":         18,
	"eare.biz":          30,
	"eayo.biz":          30,
	"eazzy.biz":         27,
	"edgeblab.biz":      29,
	"edgeclub.biz":      29,
	"edgeify.biz":       36,
	"edgepulse.biz":     21,
	"edgetag.biz":       24,
	"edgewire.biz":      29,
	"eidel.biz":         33,
	"eimbee.biz":        22,
	"einti.biz":         19,
	"eire.biz":          28,
	"fadeo.biz":         35,
	"fanoodle.biz":      23,
	"fatz.biz":          30,
	"feedbug.biz":       29,
	"feedfire.biz":      30,
	"feedfish.biz":      35,
	"feedmix.biz":       31,
	"feednation.biz":    24,
	"feedspan.biz":      28,
	"fivebridge.biz":    20,
	"fivechat.biz":      29,
	"fiveclub.biz":      23,
	"fivespan.biz":      27,
	"flashdog.biz":      20,
	"flashpoint.biz":    35,
	"flashset.biz":      30,
	"flashspan.biz":     32,
	"flipbug.biz":       27,
	"flipopia.biz":      30,
	"flipstorm.biz":     21,
	"fliptune.biz":      29,
	"gabcube.biz":       29,
	"gabspot.biz":       24,
	"gabtune.biz":       29,
	"gabtype.biz":       29,
	"gabvine.biz":       24,
	"geba.biz":          24,
	"gevee.biz":         23,
	"gigabox.biz":       28,
	"gigaclub.biz":      25,
	"gigashots.biz":     26,
	"gigazoom.biz":      29,
	"innojam.biz":       26,
	"innotype.biz":      27,
	"innoz.biz":         24,
	"izio.biz":          26,
	"jabberbean.biz":    28,
	"jabbercube.biz":    31,
	"jabbersphere.biz":  55,
	"jabberstorm.biz":   22,
	"jabbertype.biz":    27,
	"jaloo.biz":         35,
	"jamia.biz":         33,
	"janyx.biz":         33,
	"jatri.biz":         18,
	"jaxbean.biz":       28,
	"jaxnation.biz":     21,
	"jaxspan.biz":       27,
	"jaxworks.biz":      30,
	"jayo.biz":          44,
	"jazzy.biz":         32,
	"jetpulse.biz":      25,
	"jetwire.biz":       26,
	"jumpxs.biz":        29,
	"kamba.biz":         30,
	"kanoodle.biz":      19,
	"kare.biz":          30,
	"katz.biz":          62,
	"kaymbo.biz":        34,
	"kayveo.biz":        22,
	"kazio.biz":         21,
	"kazu.biz":          16,
	"kimia.biz":         25,
	"kwideo.biz":        17,
	"kwilith.biz":       25,
	"kwimbee.biz":       34,
	"kwinu.biz":         15,
	"lajo.biz":          20,
	"latz.biz":          24,
	"layo.biz":          32,
	"lazz.biz":          27,
	"lazzy.biz":         26,
	"leenti.biz":        26,
	"leexo.biz":         32,
	"linkbridge.biz":    38,
	"linkbuzz.biz":      24,
	"linklinks.biz":     31,
	"linktype.biz":      31,
	"livefish.biz":      31,
	"livepath.biz":      23,
	"livetube.biz":      53,
	"livez.biz":         28,
	"meedoo.biz":        23,
	"meejo.biz":         24,
	"meembee.biz":       26,
	"meemm.biz":         23,
	"meetz.biz":         33,
	"meevee.biz":        62,
	"meeveo.biz":        27,
	"meezzy.biz":        24,
	"miboo.biz":         26,
	"midel.biz":         28,
	"minyx.biz":         25,
	"mita.biz":          29,
	"mudo.biz":          36,
	"muxo.biz":          25,
	"mybuzz.biz":        32,
	"mycat.biz":         32,
	"mydeo.biz":         20,
	"mydo.biz":          30,
	"mymm.biz":          21,
	"mynte.biz":         54,
	"myworks.biz":       27,
	"nlounge.biz":       25,
	"npath.biz":         33,
	"ntag.biz":          28,
	"ntags.biz":         32,
	"oba.biz":           22,
	"oloo.biz":          19,
	"omba.biz":          26,
	"ooba.biz":          27,
	"oodoo.biz":         30,
	"oozz.biz":          22,
	"oyoba.biz":         27,
	"oyoloo.biz":        30,
	"oyonder.biz":       29,
	"oyondu.biz":        23,
	"oyope.biz":         24,
	"oyoyo.biz":         32,
	"ozu.biz":           18,
	"photobean.biz":     25,
	"photobug.biz":      57,
	"photofeed.biz":     25,
	"photojam.biz":      35,
	"photolist.biz":     19,
	"photospace.biz":    33,
	"pixoboo.biz":       14,
	"pixonyx.biz":       30,
	"pixope.biz":        32,
	"plajo.biz":         32,
	"plambee.biz":       29,
	"podcat.biz":        31,
	"quamba.biz":        31,
	"quatz.biz":         54,
	"quaxo.biz":         25,
	"quimba.biz":        25,
	"quimm.biz":         33,
	"quinu.biz":         60,
	"quire.biz":         25,
	"realblab.biz":      32,
	"realbridge.biz":    30,
	"realbuzz.biz":      22,
	"realcube.biz":      57,
	"realfire.biz":      37,
	"reallinks.biz":     25,
	"realmix.biz":       27,
	"realpoint.biz":     22,
	"rhybox.biz":        30,
	"rhycero.biz":       28,
	"rhyloo.biz":        32,
	"rhynoodle.biz":     25,
	"rhynyx.biz":        17,
	"rhyzio.biz":        36,
	"riffpath.biz":      21,
	"riffpedia.biz":     33,
	"riffwire.biz":      31,
	"roodel.biz":        29,
	"roombo.biz":        29,
	"roomm.biz":         32,
	"rooxo.biz":         34,
	"shufflebeat.biz":   32,
	"shuffledrive.biz":  25,
	"shufflester.biz":   26,
	"shuffletag.biz":    23,
	"skaboo.biz":        35,
	"skajo.biz":         26,
	"skalith.biz":       30,
	"skiba.biz":         22,
	"skibox.biz":        27,
	"skidoo.biz":        24,
	"skilith.biz":       29,
	"skimia.biz":        45,
	"skinder.biz":       25,
	"skinix.biz":        23,
	"skinte.biz":        39,
	"skipfire.biz":      29,
	"skippad.biz":       26,
	"skipstorm.biz":     30,
	"skiptube.biz":      26,
	"skivee.biz":        34,
	"skyba.biz":         40,
	"skyble.biz":        32,
	"skyndu.biz":        32,
	"skynoodle.biz":     28,
	"skyvu.biz":         34,
	"snaptags.biz":      33,
	"tagcat.biz":        33,
	"tagchat.biz":       37,
	"tagfeed.biz":       30,
	"tagopia.biz":       17,
	"tagpad.biz":        28,
	"tagtune.biz":       22,
	"talane.biz":        22,
	"tambee.biz":        24,
	"tanoodle.biz":      38,
	"tavu.biz":          37,
	"tazz.biz":          27,
	"tazzy.biz":         28,
	"tekfly.biz":        31,
	"teklist.biz":       26,
	"thoughtbeat.biz":   30,
	"thoughtblab.biz":   24,
	"thoughtbridge.biz": 30,
	"thoughtmix.biz":    33,
	"thoughtsphere.biz": 20,
	"thoughtstorm.biz":  38,
	"thoughtworks.biz":  24,
	"topdrive.biz":      35,
	"topicblab.biz":     32,
	"topiclounge.biz":   21,
	"topicshots.biz":    30,
	"topicstorm.biz":    22,
	"topicware.biz":     35,
	"topiczoom.biz":     38,
	"trilia.biz":        28,
	"trilith.biz":       25,
	"trudeo.biz":        29,
	"trudoo.biz":        28,
	"trunyx.biz":        33,
	"trupe.biz":         34,
	"twimbo.biz":        19,
	"twimm.biz":         30,
	"twinder.biz":       28,
	"twinte.biz":        33,
	"twitterbeat.biz":   33,
	"twitterbridge.biz": 20,
	"twitterlist.biz":   26,
	"twitternation.biz": 22,
	"twitterwire.biz":   21,
	"twitterworks.biz":  39,
	"twiyo.biz":         37,
	"vidoo.biz":         28,
	"vimbo.biz":         21,
	"vinder.biz":        31,
	"vinte.biz":         34,
	"vipe.biz":          25,
	"vitz.biz":          26,
	"viva.biz":          30,
	"voolia.biz":        34,
	"voolith.biz":       26,
	"voomm.biz":         61,
	"voonder.biz":       32,
	"voonix.biz":        32,
	"voonte.biz":        26,
	"voonyx.biz":        25,
	"wikibox.biz":       27,
	"wikido.biz":        21,
	"wikivu.biz":        23,
	"wikizz.biz":        61,
	"wordify.biz":       28,
	"wordpedia.biz":     25,
	"wordtune.biz":      27,
	"wordware.biz":      19,
	"yabox.biz":         24,
	"yacero.biz":        34,
	"yadel.biz":         27,
	"yakidoo.biz":       21,
	"yakijo.biz":        29,
	"yakitri.biz":       26,
	"yambee.biz":        20,
	"yamia.biz":         17,
	"yata.biz":          25,
	"yodel.biz":         26,
	"yodo.biz":          21,
	"yodoo.biz":         24,
	"yombu.biz":         29,
	"yotz.biz":          26,
	"youbridge.biz":     40,
	"youfeed.biz":       32,
	"youopia.biz":       22,
	"youspan.biz":       59,
	"youtags.biz":       22,
	"yoveo.biz":         31,
	"yozio.biz":         33,
	"zava.biz":          29,
	"zazio.biz":         18,
	"zoombeat.biz":      28,
	"zoombox.biz":       30,
	"zoomcast.biz":      38,
	"zoomdog.biz":       29,
	"zoomlounge.biz":    25,
	"zoomzone.biz":      32,
	"zoonder.biz":       29,
	"zoonoodle.biz":     27,
	"zooveo.biz":        22,
	"zoovu.biz":         38,
	"zooxo.biz":         33,
	"zoozzy.biz":        23,
}
