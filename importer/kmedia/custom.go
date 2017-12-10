package kmedia

import (
	"database/sql"
	"encoding/json"
	"runtime/debug"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func ImportCustom() {
	clock := Init()

	stats = NewImportStatistics()

	dump, err := loadDump()
	utils.Must(err)
	log.Infof("len(dump) = %d", len(dump))

	utils.Must(importCustomContainers(dump))

	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadDump() ([]int, error) {
	var rozaDump = []byte(`[[5958],[2755],[31273],[65981],[15357],[15667],[56116],[12445],[12446],[12722],[13554],[15144],[15145],[60535],[60567],[13553],[2018],[12214],[1746],[51759],[45026],[604],[6141],[6140],[43377],[609],[3469],[3481],[5966],[6033],[13451],[59831],[6762],[6272],[6451],[6551],[6606],[6652],[6989],[6845],[36703],[6777],[2551],[16750],[16803],[16879],[16964],[17340],[17635],[17734],[17778],[18199],[18296],[18379],[18960],[19200],[19650],[19845],[20419],[3089],[19490],[19739],[48976],[890],[8277],[54429],[698],[6236],[6137],[6209],[6182],[15358],[8943],[1001],[43761],[8880],[45287],[43759],[56115],[56905],[54793],[7610],[16413],[54788],[7611],[7612],[7613],[7614],[4476],[12215],[27064],[28173],[7615],[7616],[7617],[5956],[6031],[6142],[6208],[6196],[6271],[6359],[6399],[6450],[6552],[6605],[6653],[6763],[6846],[6934],[6983],[7036],[7085],[12787],[12870],[12932],[12913],[12973],[12974],[13043],[13044],[13108],[13109],[13165],[13166],[13222],[13223],[13283],[13284],[13342],[13452],[13453],[13499],[13618],[13619],[13739],[13736],[13802],[13803],[13870],[13871],[13939],[14000],[14001],[14071],[14072],[14139],[14140],[14186],[14391],[14464],[14525],[14526],[14592],[14593],[14787],[14836],[14837],[14907],[14974],[14975],[15062],[27052],[7618],[38905],[56855],[7619],[43375],[7596],[7621],[9020],[9065],[7622],[9273],[9705],[9771],[7623],[5959],[6143],[6235],[6136],[15209],[15210],[15284],[15359],[15360],[15441],[15527],[15528],[15593],[15594],[15668],[15669],[15794],[15795],[15879],[15880],[16127],[20307],[20312],[45285],[535],[539],[541],[572],[578],[587],[700],[702],[706],[709],[801],[816],[12257],[15962],[16005],[16486],[16590],[16591],[16660],[16661],[16751],[16804],[16877],[16878],[16965],[17049],[17125],[17126],[17341],[17342],[17425],[17515],[17516],[17636],[17735],[17779],[17780],[17876],[17877],[18019],[18200],[18201],[18298],[18378],[18703],[18859],[18961],[18962],[19201],[19488],[19572],[19651],[19652],[19737],[19738],[19846],[20188],[59830],[6688],[48542],[54796],[56902],[17264],[17263],[65980],[13163],[8668],[6237],[14591],[16485],[18202],[15595],[19653],[56903],[16662],[13800],[14137],[17736],[13620],[18963],[54795],[54807],[14909],[13873],[12788],[12258],[20308],[13220],[17426],[13221],[14185],[43374],[13621],[19818],[13801],[12259],[14976],[60533],[14590],[7532],[54785],[14528],[12789],[15146],[6453],[28093],[16663],[13500],[19348],[56851],[15442],[20812],[9012],[13940],[8817],[15674],[19202],[17050],[18702],[13735],[16880],[15361],[14390],[18860],[15063],[14838],[16128],[15211],[16129],[16414],[16805],[16752],[18020],[49925],[13449],[14138],[13341],[19740],[14977],[14910],[60565],[12723],[12724],[16592],[6291],[15285],[6292],[6358],[6400],[6553],[6554],[6607],[6654],[6764],[6847],[6935],[6936],[6984],[7038],[7037],[7087],[7088],[7139],[7182],[7232],[54427],[7273],[17878],[17781],[14002],[17637],[15529],[8099],[13042],[15212],[15147],[11289,56023],[13555],[14788],[9022],[14731],[13110],[14463],[14004],[18299],[12447],[14527],[13872],[56900],[19741],[54794],[13167],[56856],[19221],[19222],[19223],[19224],[19225],[19226],[19227],[19228],[19229],[19230],[20945],[19231],[19232],[19233],[19234],[19235],[19236],[19237],[19238],[19239],[19240],[19605],[19241],[19242],[19243],[19244],[19245],[19248],[19249],[19250],[19251],[19246],[19252],[19253],[19254],[19255],[19256],[19257],[19258],[19259],[19260],[22860],[19265],[19261],[19262],[19263],[19264],[19266],[19267],[19268],[19269],[19270],[19271],[22861],[19272],[19276],[19273],[19274],[19275],[19277],[19278],[19279],[19280],[19281],[19282],[19283],[19284],[19285],[19286],[19287],[19288],[19289],[19290],[19291],[19292],[7884],[8382],[9835],[14943],[14944],[14945],[15074],[15075],[7868],[7869],[7870],[8019],[7871],[7872],[7873],[7874],[7875],[7876],[7877],[7878],[7879],[8368],[8369],[7879],[8411],[14938],[14939],[14940],[14941],[15072],[7880],[7881],[7882],[9834],[14942],[9834],[47177],[7883],[25485],[6709],[65470,65471,65472],[5123],[5130],[5957],[6034],[6765],[6937],[7089],[8882],[9064],[9093],[11856],[12206],[12256],[12264],[12444],[12704],[12725],[12866],[12911],[13040],[13105],[13090],[13164],[13217],[13281],[13337],[13445],[13498],[13557],[13617],[13733],[13797],[13868],[13936],[13999],[14066],[14136],[14183],[14388],[14460],[14523],[14587],[14730],[14789],[14841],[14908],[14979],[15061],[15149],[15214],[15282],[15362],[15444],[15531],[15597],[15675],[15798],[15877],[15878],[15961],[16007],[16131],[16132],[16664],[16753],[16806],[16881],[16882],[16963],[17048],[17265],[17428],[17740],[17782],[17875],[18018],[18260],[18197],[18704],[18858],[18958],[18959],[19199],[19487],[19654],[20190],[20418],[20309],[20313],[22519],[19187],[19576],[19773],[19693],[20121],[20175],[21194],[21755],[23845],[24403],[27430],[28285],[53375],[48523],[49551],[57788],[129],[23274],[34713],[54337],[38943],[2617],[2826],[2825],[3178],[6138],[6167,9175],[6938],[7355],[7395],[7406],[9097],[9175],[10146],[12871],[12971],[13282],[13030],[16487],[16783],[16783],[16847],[17214],[18203],[24302],[24553],[24554],[24917],[24900],[24901],[25410],[25411],[25397],[27501],[25719],[25720],[26100],[26099],[27203],[27031],[27034],[27468,27469],[27760],[28515],[32070],[33942],[33947],[33946],[35719],[35724],[35728],[35729],[35732],[35726],[35734],[39036,39034],[38752,38767],[38754,38768],[38765,38769],[38770,38771],[38793],[38792],[42165,42166],[42168,42176],[42170,42171],[42173,42174],[42113,42112],[43177,43172],[43175,43176],[44402],[48177],[48175],[48173],[48188],[48185],[47322],[47746],[49262,49260,49265,49263,49264,49254,49255,49256,49258,49257,49261,49259],[54357],[54361],[54363],[54364],[54368],[54374],[54377],[54358],[54362],[54367],[54369],[54373],[54376],[56064],[56071],[56074],[56060],[56062],[56068],[63281,63280,63279,63275,63278,63276,63277],[59829],[64624],[64628],[64629],[67041],[12328],[12330],[12331],[12332],[12333],[14732],[14733],[14734],[14735],[14736],[14738],[14740],[19571],[20503],[27053],[27065],[7826],[7837],[7841],[7846],[7839],[7856,14609],[7847],[2296],[7853],[7843],[7836],[7834],[7835],[14357],[7842],[14356],[7838],[7840],[7851],[7855],[13285],[13340],[13339],[13448],[13447],[13446],[13622],[13734],[13799],[13869],[13938],[14389],[14462],[14589],[34613],[35346],[35622],[36117],[36521],[37268,37266],[38448,38449],[39501,39502],[40803],[44609],[44735,44736],[45362,45349],[46829],[47217],[48000],[48555,48462],[27292],[31753],[34844],[35201],[35350],[35893],[35928],[35983],[36005],[36163],[36758],[36939,36948],[37058],[37348,37357],[37521,37526],[37516,37520],[37868,38533,38750],[37869,37879],[38085],[38129,38130],[38177,38178],[38269],[38297],[38340,38367],[38582],[38681,38683],[38900,38901],[39214,39061],[39351,39451],[39647,39648],[39833,39834],[39875,39876],[40083,40084],[41577,41578],[41633,41632],[44329,44331],[44399],[44943,44945],[45033,45004],[45199],[46032],[45973],[46233,46231],[46345],[46346],[46292],[47422],[47561],[47784],[47921],[48385],[48462],[6195],[6153],[6154],[6156],[6183],[6475],[6564],[6576],[6767],[6829],[6848],[7020],[7092],[7067],[7140],[7120],[7201],[7167],[7211],[7279],[7338],[7331],[7392],[7422],[7467],[7535],[7563],[8070],[8069],[7940],[8005],[8507],[8816],[15992],[17853],[19365],[19548],[19588],[20447],[20448],[14524],[26493],[5540],[2554],[2708],[2765],[2925],[2923],[2980],[4714],[5985],[6234],[6289],[6561],[6562],[6563],[6849],[6939],[6985],[6999],[6987],[7039],[7091],[7145],[7141],[7184],[7185],[7354],[7798],[7926],[7960],[8006],[8007],[8056],[8772],[8871],[8887],[8913],[12213],[12728],[12912],[13218],[13558],[13559],[13796],[14387],[14386],[14459],[14842],[15286],[15363],[15676],[15677],[15799],[15965],[16009],[16154],[16412],[16488],[16665],[17124],[17262],[17345],[17429],[17518],[17517],[17519],[17638],[17639],[17640],[17738],[17739],[17783],[17784],[17879],[17880],[18198],[19589],[19819],[21023],[26205],[26494],[26819],[27588],[27587],[36723],[40195,40161],[40199,40200],[40279,40278],[3957],[25305],[18705],[18705],[20500],[20923],[21567],[7611],[7613],[7614],[7616]]`)
	var data [][]int
	if err := json.Unmarshal(rozaDump, &data); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	cnIDs := make([]int, len(data))
	for i := range data {
		if len(data[i]) == 1 {
			cnIDs[i] = data[i][0]
		} else {
			cnIDs[i] = -1
		}
	}

	return cnIDs, nil
}

type ImportSettings struct {
	CT             string
	SubCT          string
	CollectionUIDs []string
	KmCNID         int
}

func getImportSettings(idx int) *ImportSettings {
	i := idx
	if idx > 395 {
		i--
	}
	switch {
	case 1 <= i && i <= 394,
		510 <= i && i <= 599,
		723 <= i && i <= 738,
		i == 746,
		759 <= i && i <= 772,
		790 <= i && i <= 840,
		971 <= i && i <= 975:
		return &ImportSettings{CT: api.CT_LESSON_PART}
	case i == 600:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"o06JwHkK"}}
	case i == 601:
		return &ImportSettings{CT: api.CT_CLIP, CollectionUIDs: []string{"iNQzqlKk"}}
	case i == 602:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"iNQzqlKk"}}
	case i == 603:
		return &ImportSettings{CT: api.CT_FRIENDS_GATHERING, CollectionUIDs: []string{"iNQzqlKk"}}
	case i == 604:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"tns8a19k"}}
	case i == 605:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"48mAodgk"}}
	case i == 606:
		return &ImportSettings{CT: api.CT_CLIP, CollectionUIDs: []string{"LRD8aEXh"}}
	case i == 607:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"y11bttvS"}}
	case i == 608:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"jRGTUnhf"}}
	case i == 609:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"dONrFRuT"}}
	case i == 610:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"PSyY0wcr"}}
	case i == 611:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"ShcTNcKD"}}
	case 612 <= i && i <= 614:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"egGmzRy6"}}
	case i == 615:
		return &ImportSettings{CT: api.CT_FRIENDS_GATHERING, CollectionUIDs: []string{"jDqE3wjq"}}
	case 616 <= i && i <= 619:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"egGmzRy6"}}
	case i == 620:
		return &ImportSettings{CT: api.CT_CLIP}
	case i == 621:
		return &ImportSettings{CT: api.CT_FRIENDS_GATHERING, CollectionUIDs: []string{"5KNIGkeP"}}
	case i == 624:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"TGPflh11"}}
	case i == 625:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"KtD5ALuS"}}
	case 627 <= i && i <= 630:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"3dNYLvrq"}}
	case i == 631:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"E66uiftK"}}
	case i == 633:
		return &ImportSettings{CT: api.CT_CLIP, CollectionUIDs: []string{"VSmQzPFg"}}
	case i == 634:
		return &ImportSettings{CT: api.CT_MEAL, CollectionUIDs: []string{"16tQIHD7"}}
	case 635 <= i && i <= 636:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"CGCZpSC8"}}
	case 638 <= i && i <= 641:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"phnhxbPu"}}
	case i == 642:
		return &ImportSettings{CT: api.CT_EVENT_PART, SubCT: "EREV_TARBUT", CollectionUIDs: []string{"phnhxbPu"}}
	case i == 643:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"ooOga3GM"}}
	case i == 644:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"WWxzJL4R"}}
	case i == 645:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"km84LtML"}}
	case i == 646:
		return &ImportSettings{CT: api.CT_EVENT_PART, SubCT: "EVENT", CollectionUIDs: []string{"km84LtML"}}
	case i == 647:
		return &ImportSettings{CT: api.CT_EVENT_PART, SubCT: "EVENT", CollectionUIDs: []string{"off8IvGD"}}
	case i == 648:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"off8IvGD"}}
	case i == 649:
		return &ImportSettings{CT: api.CT_EVENT_PART, SubCT: "EVENT", CollectionUIDs: []string{"off8IvGD"}}
	case i == 650:
		return &ImportSettings{CT: api.CT_EVENT_PART, SubCT: "EVENT", CollectionUIDs: []string{"ks245xT9"}}
	case i == 651:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"off8IvGD"}}
	case i == 652:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"pKT8Ytcd"}}
	case i == 654:
		return &ImportSettings{CT: api.CT_EVENT_PART, SubCT: "EVENT", CollectionUIDs: []string{"98XMWsYU"}}
	case i == 655:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"98XMWsYU"}}
	case i == 656:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"WUexJyuq"}}
	case i == 657:
		return &ImportSettings{CT: api.CT_EVENT_PART, SubCT: "EVENT", CollectionUIDs: []string{"WUexJyuq"}}
	case i == 659:
		return &ImportSettings{CT: api.CT_EVENT_PART, SubCT: "EVENT", CollectionUIDs: []string{"g0CSoUru", "TDu5WjB7"}}
	case i == 660:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"g0CSoUru", "TDu5WjB7"}}
	case i == 661:
		return &ImportSettings{CT: api.CT_EVENT_PART, SubCT: "EVENT", CollectionUIDs: []string{"gQwX8atP", "TDu5WjB7"}}
	case i == 663:
		return &ImportSettings{CT: api.CT_LECTURE, CollectionUIDs: []string{"mv36UoEE"}}
	case i == 664:
		return &ImportSettings{CT: api.CT_VIDEO_PROGRAM_CHAPTER, CollectionUIDs: []string{"KkZ63W0r", "YvPmFc3Z"}}
	case 665 <= i && i <= 667:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"gHf0tWkr"}}
	case 668 <= i && i <= 672:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"Cc36c1Cj"}}
	case 673 <= i && i <= 674:
		return &ImportSettings{CT: api.CT_FRIENDS_GATHERING, CollectionUIDs: []string{"Cc36c1Cj"}}
	case i == 675:
		return &ImportSettings{CT: api.CT_WOMEN_LESSON, CollectionUIDs: []string{"sUA9vBdk"}}
	case 676 <= i && i <= 679:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"IfrZxKQV"}}
	case 680 <= i && i <= 681:
		return &ImportSettings{CT: api.CT_FRIENDS_GATHERING, CollectionUIDs: []string{"IfrZxKQV"}}
	case 682 <= i && i <= 685:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"jVMINyzi"}}
	case i == 686:
		return &ImportSettings{CT: api.CT_WOMEN_LESSON, CollectionUIDs: []string{"jVMINyzi"}}
	case 687 <= i && i <= 688:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"dC3GuoJO"}}
	case i == 689:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"uH1rc4V8"}}
	case i == 690:
		return &ImportSettings{CT: api.CT_MEAL, CollectionUIDs: []string{"4TK3vbwB"}}
	case 691 <= i && i <= 694:
		return &ImportSettings{CT: api.CT_TRAINING, CollectionUIDs: []string{"4TK3vbwB"}}
	case 695 <= i && i <= 696:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"4TK3vbwB"}}
	case i == 697:
		return &ImportSettings{CT: api.CT_CLIP, CollectionUIDs: []string{"4TK3vbwB"}}
	case 698 <= i && i <= 704:
		return &ImportSettings{CT: api.CT_FRIENDS_GATHERING, CollectionUIDs: []string{"NvzsYWBk"}}
	case 705 <= i && i <= 710:
		return &ImportSettings{CT: api.CT_MEAL, CollectionUIDs: []string{"NvzsYWBk"}}
	case 711 <= i && i <= 713:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"PMVJaXHs"}}
	case 714 <= i && i <= 716:
		return &ImportSettings{CT: api.CT_FRIENDS_GATHERING, CollectionUIDs: []string{"PMVJaXHs"}}
	case i == 717:
		return &ImportSettings{CT: api.CT_CLIP, CollectionUIDs: []string{"YKXIJQQ5"}}
	case i == 718:
		return &ImportSettings{CT: api.CT_LESSON_PART, CollectionUIDs: []string{"sAeBdIME", "o8QveYwt"}}
	case 719 <= i && i <= 721:
		return &ImportSettings{CT: api.CT_FRIENDS_GATHERING, CollectionUIDs: []string{"ujshviDc"}}
	case i == 722:
		return &ImportSettings{CT: api.CT_MEAL, CollectionUIDs: []string{"Q7gbBk7g"}}
	case 773 <= i && i <= 788:
		return &ImportSettings{CT: api.CT_WOMEN_LESSON}
	case i == 789:
		return &ImportSettings{CT: api.CT_VIDEO_PROGRAM_CHAPTER, CollectionUIDs: []string{"ZJZVhI9z"}}
	case 841 <= i && i <= 880:
		return &ImportSettings{CT: api.CT_MEAL}
	case i == 970:
		return &ImportSettings{CT: api.CT_MEAL}
	case i == 976:
		return &ImportSettings{CT: api.CT_VIDEO_PROGRAM_CHAPTER, CollectionUIDs: []string{"jSPXwGwQ"}}
	default:
		return nil
	}
}

func importCustomContainers(dump []int) error {
	noCN := make([]int, 0)
	noSettings := make([]int, 0)

	for i := range dump {
		idx := i + 1
		if 395 <= idx && idx <= 506 {
			continue
		}

		settings := getImportSettings(idx)
		if settings == nil {
			noSettings = append(noSettings, idx)
			continue
		}

		if dump[i] <= 0 {
			noCN = append(noCN, idx)
			continue
		}

		settings.KmCNID = dump[i]

		tx, err := mdb.Begin()
		utils.Must(err)

		if err = importNewCUBySettings(tx, settings); err != nil {
			utils.Must(tx.Rollback())
			stats.TxRolledBack.Inc(1)
			log.Error(err)
			debug.PrintStack()
			continue
		} else {
			utils.Must(tx.Commit())
			stats.TxCommitted.Inc(1)
		}
	}

	log.Infof("len(noSettings) = %d", len(noSettings))
	for i := range noSettings {
		log.Infof("%d", noSettings[i])
	}

	log.Infof("len(noCN) = %d", len(noCN))
	for i := range noCN {
		log.Infof("%d", noCN[i])
	}

	return nil
}

func importNewCUBySettings(exec boil.Executor, settings *ImportSettings) error {
	cu, err := models.ContentUnits(mdb, qm.Where("(properties->>'kmedia_id')::int = ?", settings.KmCNID)).One()
	if err != nil {
		if err != sql.ErrNoRows {
			return errors.Wrapf(err, "Lookup content_unit kmedia_id=%d", settings.KmCNID)
		}
	} else {
		return errors.Errorf("CU exists KmCNID=%d cu.ID %d. skipping", settings.KmCNID, cu.ID)
	}

	cn, err := kmodels.FindContainer(kmdb, settings.KmCNID)
	if err != nil {
		return errors.Wrapf(err, "Lookup container %d", settings.KmCNID)
	}

	cs := make([]*models.Collection, 0)
	if settings.CollectionUIDs != nil {
		for j := range settings.CollectionUIDs {
			uid := settings.CollectionUIDs[j]
			c, err := models.Collections(exec, qm.Where("uid = ?", uid)).One()
			if err != nil {
				return errors.Wrapf(err, "Lookup collection uid = %s", uid)
			}
			cs = append(cs, c)
		}
	}

	// create unit
	cu, err = importContainerWOCollectionNewCU(exec, cn, settings.CT)
	if err != nil {
		return errors.Wrapf(err, "Import new container %d", cn.ID)
	}

	// create or CCUs
	for i := range cs {
		c := cs[i]

		log.Infof("Associating %d %s to %s: [cu,c]=[%d,%d]", cn.ID, cn.Name.String, c.UID, cu.ID, c.ID)

		ccuName := strconv.Itoa(cn.Position.Int)
		if settings.CT == api.CT_EVENT_PART {
			ccuName = settings.SubCT + ccuName
		}

		err = createOrUpdateCCU(exec, cu, models.CollectionsContentUnit{
			CollectionID:  c.ID,
			ContentUnitID: cu.ID,
			Name:          ccuName,
			Position:      cn.Position.Int,
		})
		if err != nil {
			return errors.Wrapf(err, "Create or update CCU %d", cn.ID)
		}

		if !c.Published && cu.Published {
			c.Published = true
			if err := c.Update(exec); err != nil {
				return errors.Wrapf(err, "Update collection.published %d", c.ID)
			}
		}
	}

	return nil
}
