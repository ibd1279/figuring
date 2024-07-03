package figuring

// see https://pomax.github.io/bezierinfo/legendre-gauss.html
var (
	legendregauss_abscissa [64]float64 = [64]float64{
		-0.0243502926634244325089558428537156614268871093149758091634531663960566965166295288529853061657116894882370493013671717560479926679408068852617342586968190919443025679363843727751902756254975073084367002129407854253246662805532069172532219089321005870178809284335033318073251039701073379759,
		0.0243502926634244325089558428537156614268871093149758091634531663960566965166295288529853061657116894882370493013671717560479926679408068852617342586968190919443025679363843727751902756254975073084367002129407854253246662805532069172532219089321005870178809284335033318073251039701073379759,
		-0.0729931217877990394495429419403374932441261816687788533563163323377395217254429050181833064967505478802134768007678458612956459126148837496307967995621683597067400860540057918571357609346700883624064782909888895547912499697295335516804810990011835717296819569741981551097569810739977931249,
		0.0729931217877990394495429419403374932441261816687788533563163323377395217254429050181833064967505478802134768007678458612956459126148837496307967995621683597067400860540057918571357609346700883624064782909888895547912499697295335516804810990011835717296819569741981551097569810739977931249,
		-0.1214628192961205544703764634922478782186836383371912940423495826006931832245070944213952236889690237679661122848626352437115925113582515415979746598665681268376919737373113667247142315234607397986222184384307059013512155412722263090766858403726100735651684437098923088469949297570988091677,
		0.1214628192961205544703764634922478782186836383371912940423495826006931832245070944213952236889690237679661122848626352437115925113582515415979746598665681268376919737373113667247142315234607397986222184384307059013512155412722263090766858403726100735651684437098923088469949297570988091677,
		-0.1696444204239928180373136297482698441999902667343778505894178746884342357796505591325925801106834127602396624627746208498838190598644711533868179088757129652678801285453336384132061358206434768209251779433993367981112053126660575920785488821886662635718276505127786854167528795165050389903,
		0.1696444204239928180373136297482698441999902667343778505894178746884342357796505591325925801106834127602396624627746208498838190598644711533868179088757129652678801285453336384132061358206434768209251779433993367981112053126660575920785488821886662635718276505127786854167528795165050389903,
		-0.217423643740007084149648748988822617578485831141222348630380401885689634737659235537163737740243604800759921292790013642836998201691226098978544332296437547642195961469059807833597096166848933098833151166287901339013986496737408125314858259377210847115061960167857239951500335854820530586,
		0.217423643740007084149648748988822617578485831141222348630380401885689634737659235537163737740243604800759921292790013642836998201691226098978544332296437547642195961469059807833597096166848933098833151166287901339013986496737408125314858259377210847115061960167857239951500335854820530586,
		-0.2646871622087674163739641725100201179804131362950932439559895448126206429452852982016450901649805445999078728714943692622330016257776575354105370883948495294882935681426154386660616476411740312465060150591301869544672050088454083442813632094277160007745547301572849406353682760306061929911,
		0.2646871622087674163739641725100201179804131362950932439559895448126206429452852982016450901649805445999078728714943692622330016257776575354105370883948495294882935681426154386660616476411740312465060150591301869544672050088454083442813632094277160007745547301572849406353682760306061929911,
		-0.3113228719902109561575126985601568835577153578680501269954571709858169098868398268719654999460149709757804582988077747605532896065842023340674450299515989484487746153929299031759475919924980452933324186984188982046762542556035347023744560814177013801414889023264804693830155588690576492164,
		0.3113228719902109561575126985601568835577153578680501269954571709858169098868398268719654999460149709757804582988077747605532896065842023340674450299515989484487746153929299031759475919924980452933324186984188982046762542556035347023744560814177013801414889023264804693830155588690576492164,
		-0.357220158337668115950442615046202531626264464640909112021237019340099177403802509741325589540743874845093675632547750287037622834793938695456980400958079292460482315821150714268593539935795231095157602025909339384694681190969656053235824652679875951093689200190014853543993102190088381483,
		0.357220158337668115950442615046202531626264464640909112021237019340099177403802509741325589540743874845093675632547750287037622834793938695456980400958079292460482315821150714268593539935795231095157602025909339384694681190969656053235824652679875951093689200190014853543993102190088381483,
		-0.4022701579639916036957667712601588487132652056150208082760431843129087214967261515969669708970990221669508217089555714806012046537594438323569293594638517933840725639831594134038262580440842200076281605641993773325072728778928440394419613403725280285705765326861533477990551765453978736181,
		0.4022701579639916036957667712601588487132652056150208082760431843129087214967261515969669708970990221669508217089555714806012046537594438323569293594638517933840725639831594134038262580440842200076281605641993773325072728778928440394419613403725280285705765326861533477990551765453978736181,
		-0.4463660172534640879849477147589151892067507578262501763082606820212937626970791295735937384813941473610238854736863966831464694923749954564921955859791688348936235671762050333408576202492209167272366373825152067680845198006626563761196191045700093968519269790165913301841545609485952718504,
		0.4463660172534640879849477147589151892067507578262501763082606820212937626970791295735937384813941473610238854736863966831464694923749954564921955859791688348936235671762050333408576202492209167272366373825152067680845198006626563761196191045700093968519269790165913301841545609485952718504,
		-0.4894031457070529574785263070219213908493732974637398373316540793240315585537584844752851087115581833443158831657759501916744927211636581386025171070582998790865902140901838128045602667002106847665927788023098753138400106615804847725751247952878407027140260429761863258319891988431055536872,
		0.4894031457070529574785263070219213908493732974637398373316540793240315585537584844752851087115581833443158831657759501916744927211636581386025171070582998790865902140901838128045602667002106847665927788023098753138400106615804847725751247952878407027140260429761863258319891988431055536872,
		-0.531279464019894545658013903544455247408525588734180238053268047166041778398245121448843253296460411619816073385211875151397248937264089998182375345915413219579220233566173902955487674957069948591213673456625912506280248298229907928060620469290406581396192419570799497688513058132396498814,
		0.531279464019894545658013903544455247408525588734180238053268047166041778398245121448843253296460411619816073385211875151397248937264089998182375345915413219579220233566173902955487674957069948591213673456625912506280248298229907928060620469290406581396192419570799497688513058132396498814,
		-0.5718956462026340342838781166591886431831910060912509932273284719418912212643223327597417735844776972648163821225207266145263395898251858906124801356522395225326954546582593857056725545247314092886133249957455688586118199388064447508712958637376347457406936802691416300157802889354368128467,
		0.5718956462026340342838781166591886431831910060912509932273284719418912212643223327597417735844776972648163821225207266145263395898251858906124801356522395225326954546582593857056725545247314092886133249957455688586118199388064447508712958637376347457406936802691416300157802889354368128467,
		-0.6111553551723932502488529710185489186961245593079718443367976666933088374650288148448667879830703867726577720666491772560110368248450475818132595468834493579434418252468978282181164008820769765174450056817468275783966537351796625224747700783378315186174632657840114512887287754747902924865,
		0.6111553551723932502488529710185489186961245593079718443367976666933088374650288148448667879830703867726577720666491772560110368248450475818132595468834493579434418252468978282181164008820769765174450056817468275783966537351796625224747700783378315186174632657840114512887287754747902924865,
		-0.6489654712546573398577612319934048855296904334732011728792502624366057598865738239773166627826358871699142531853930525866830933399401844502541092962631127742267449897116125748014270680393024011359139031312062573520509858894743036198986443014969748157931850178889912972291202354657382925509,
		0.6489654712546573398577612319934048855296904334732011728792502624366057598865738239773166627826358871699142531853930525866830933399401844502541092962631127742267449897116125748014270680393024011359139031312062573520509858894743036198986443014969748157931850178889912972291202354657382925509,
		-0.6852363130542332425635583710313763019356410785396718681324042749913611976548967332647625541234155413035075739348863233240851709341392873173633850612006618690218164007761541855237208605116527909791956398350719463021018362527198358286721239529091637248252469435642287693207506339068528700205,
		0.6852363130542332425635583710313763019356410785396718681324042749913611976548967332647625541234155413035075739348863233240851709341392873173633850612006618690218164007761541855237208605116527909791956398350719463021018362527198358286721239529091637248252469435642287693207506339068528700205,
		-0.7198818501716108268489402178319472447581380033149019526220473151184468592486433646042300919262498902882179891494724046747921544602557246455427317830132976771174504209146835854012577372395960646854355024460204286901035470812684876587693044068989973704915078171158689213412715251485203442313,
		0.7198818501716108268489402178319472447581380033149019526220473151184468592486433646042300919262498902882179891494724046747921544602557246455427317830132976771174504209146835854012577372395960646854355024460204286901035470812684876587693044068989973704915078171158689213412715251485203442313,
		-0.7528199072605318966118637748856939855517142713220871932461987761167722639636968539390413583009467924995905147153923286347693864945784890119950315095740891764880991728646942920208355501445550654259850385377192973526980795946453961077798744753892199929235882097232836682421729944934586281945,
		0.7528199072605318966118637748856939855517142713220871932461987761167722639636968539390413583009467924995905147153923286347693864945784890119950315095740891764880991728646942920208355501445550654259850385377192973526980795946453961077798744753892199929235882097232836682421729944934586281945,
		-0.7839723589433414076102205252137682840564141249898259334132759617476816578705098509357116190608002325895348207611752987335385494893726027026038902508685496606160441965948835252795524014713290879877269643684102005605450140247125536147801312017065532602835003540212221564314236937509990182173,
		0.7839723589433414076102205252137682840564141249898259334132759617476816578705098509357116190608002325895348207611752987335385494893726027026038902508685496606160441965948835252795524014713290879877269643684102005605450140247125536147801312017065532602835003540212221564314236937509990182173,
		-0.8132653151227975597419233380863033406981418225655972166956485149356586346082019870309280128411412936411423614767918756843380999442447282903502051218203273573634847203121451086379808399639198510674436238195505371716160648058477202993836014352158139813219612968106205248494087577632573534973,
		0.8132653151227975597419233380863033406981418225655972166956485149356586346082019870309280128411412936411423614767918756843380999442447282903502051218203273573634847203121451086379808399639198510674436238195505371716160648058477202993836014352158139813219612968106205248494087577632573534973,
		-0.840629296252580362751691544695873302982489823801755353928202683075593465893922171840726147868117503717663799561956411215937924134571068943700343442753760445948626598735504632170407243376224222403038093781056024445977626666740664628412660960413062370047183186652885532589557452614451434048,
		0.840629296252580362751691544695873302982489823801755353928202683075593465893922171840726147868117503717663799561956411215937924134571068943700343442753760445948626598735504632170407243376224222403038093781056024445977626666740664628412660960413062370047183186652885532589557452614451434048,
		-0.8659993981540928197607833850701575024125019187582496425664279511808356713122188567857456842034906362573453815878913951040194915987006979015304835979058725276345799813088989383312475641092775164460639450521468294104011206574786429237252678172922104036725327539940502197291939132802457917836,
		0.8659993981540928197607833850701575024125019187582496425664279511808356713122188567857456842034906362573453815878913951040194915987006979015304835979058725276345799813088989383312475641092775164460639450521468294104011206574786429237252678172922104036725327539940502197291939132802457917836,
		-0.8893154459951141058534040382728516224291944615104521893194744566084811090577722526400445910623711480590529533188832105988657269430913287263821624762648137092066620632787986348052306840101775313644572400860845559833367997001666659907951051347410546710134120144598833115095140475669485797579,
		0.8893154459951141058534040382728516224291944615104521893194744566084811090577722526400445910623711480590529533188832105988657269430913287263821624762648137092066620632787986348052306840101775313644572400860845559833367997001666659907951051347410546710134120144598833115095140475669485797579,
		-0.9105221370785028057563806680083298610134880848883640292531723714467102234556291968179018775780308458024302103848451312741663820589200520720207891653884985710130867134073520525932445557074805974235006810370309087879564826639263972805682465506594098949560288847385983395160311034445386606259,
		0.9105221370785028057563806680083298610134880848883640292531723714467102234556291968179018775780308458024302103848451312741663820589200520720207891653884985710130867134073520525932445557074805974235006810370309087879564826639263972805682465506594098949560288847385983395160311034445386606259,
		-0.9295691721319395758214901545592256073474270144297154975928116833612430986265594515998834499355844736686512805129688214992047597092114291955925880175797899765980745854426738149516325837607227287481909072315347776012991222301207304052068204069335766550173941103055407746774520789612843561385,
		0.9295691721319395758214901545592256073474270144297154975928116833612430986265594515998834499355844736686512805129688214992047597092114291955925880175797899765980745854426738149516325837607227287481909072315347776012991222301207304052068204069335766550173941103055407746774520789612843561385,
		-0.9464113748584028160624814913472647952793949717952331902317789712973664402149436591260928179188420533516264142755452159723722786167537167514691534968355366202934342465086995943893699962972237343218079763936958985487264411542890941861254842843026890160131607678957282346112697993618567018237,
		0.9464113748584028160624814913472647952793949717952331902317789712973664402149436591260928179188420533516264142755452159723722786167537167514691534968355366202934342465086995943893699962972237343218079763936958985487264411542890941861254842843026890160131607678957282346112697993618567018237,
		-0.9610087996520537189186141218971572067621146110378459494461586158623919945488992563976780806866203786216001498788310714552847469661399216303755820947005848739467276644122915754949838610353627723679982220628115164983443994552616161584523205789167087822341423097206088828267065770404672828066,
		0.9610087996520537189186141218971572067621146110378459494461586158623919945488992563976780806866203786216001498788310714552847469661399216303755820947005848739467276644122915754949838610353627723679982220628115164983443994552616161584523205789167087822341423097206088828267065770404672828066,
		-0.9733268277899109637418535073522726680261452944551741758819139781978152256958453749994966038154125547612207903105020176172420237675899907788807087542221018040460410464083361271842759039530092449625891215101984663282728542290395875313124045226564547294745437773482395329023327909760431499638,
		0.9733268277899109637418535073522726680261452944551741758819139781978152256958453749994966038154125547612207903105020176172420237675899907788807087542221018040460410464083361271842759039530092449625891215101984663282728542290395875313124045226564547294745437773482395329023327909760431499638,
		-0.9833362538846259569312993021568311169452475066237403837464872131233426128415470535606559721330818003585532628124845662897410684694651251174207713020897795837892725294581710205598344576799985346970638130204876060998657059283079767876980544166132523941283823202290746667358872631036031924711,
		0.9833362538846259569312993021568311169452475066237403837464872131233426128415470535606559721330818003585532628124845662897410684694651251174207713020897795837892725294581710205598344576799985346970638130204876060998657059283079767876980544166132523941283823202290746667358872631036031924711,
		-0.9910133714767443207393823834433031136413494453907904852225427459378131658644129997345108950133770434340330151289100150097018332483423277136039914249575686591612502752158650205954671083696496347591169012794322303027309768195334920157669446268175983954105533989275308193580349506657360682085,
		0.9910133714767443207393823834433031136413494453907904852225427459378131658644129997345108950133770434340330151289100150097018332483423277136039914249575686591612502752158650205954671083696496347591169012794322303027309768195334920157669446268175983954105533989275308193580349506657360682085,
		-0.9963401167719552793469245006763991232098575063402266121352522199507030568202208530946066801021703916301511794658310735397567341036554686814952726523955953805437164277655915410358813984246580862850974195805395101678543649116458555272523253307828290553873260588621490898443701779725568118502,
		0.9963401167719552793469245006763991232098575063402266121352522199507030568202208530946066801021703916301511794658310735397567341036554686814952726523955953805437164277655915410358813984246580862850974195805395101678543649116458555272523253307828290553873260588621490898443701779725568118502,
		-0.9993050417357721394569056243456363119697121916756087760628072954617646543505331997843242376462639434945376776512170265314011232493020401570891594274831367800115383317335285468800574240152992751785027563437707875403545865305271045717258142571193695943317890367167086616955235477529427992282,
		0.9993050417357721394569056243456363119697121916756087760628072954617646543505331997843242376462639434945376776512170265314011232493020401570891594274831367800115383317335285468800574240152992751785027563437707875403545865305271045717258142571193695943317890367167086616955235477529427992282,
	}
	legendregauss_weight [64]float64 = [64]float64{
		0.0486909570091397203833653907347499124426286922838743305086688042456914190998246107310291565645676057401607079939845156005172257043376703767287395573765236401039685866479381075274920900511719320271157129622463682509122641788910270632229694394595885032921037399187298767076084601033342936131,
		0.0486909570091397203833653907347499124426286922838743305086688042456914190998246107310291565645676057401607079939845156005172257043376703767287395573765236401039685866479381075274920900511719320271157129622463682509122641788910270632229694394595885032921037399187298767076084601033342936131,
		0.0485754674415034269347990667839781136875565447049294857111516761025158193093697039229163427633930410186232149083923688162761488505704450697417589593116703853157329164894580165517236877241308351214870169600093357854651930986960906313726182992933363325614247750209880050786299287510692780499,
		0.0485754674415034269347990667839781136875565447049294857111516761025158193093697039229163427633930410186232149083923688162761488505704450697417589593116703853157329164894580165517236877241308351214870169600093357854651930986960906313726182992933363325614247750209880050786299287510692780499,
		0.048344762234802957169769527158017809703692550609501080629442201445249828946429202156764153264348308119169811534137999799779908820312744765416129733427088646813066886130539178187597540312913636916139844188190193872629488730769015964208394624398401975043997268903006190530430762197842013971,
		0.048344762234802957169769527158017809703692550609501080629442201445249828946429202156764153264348308119169811534137999799779908820312744765416129733427088646813066886130539178187597540312913636916139844188190193872629488730769015964208394624398401975043997268903006190530430762197842013971,
		0.0479993885964583077281261798713460699543167134714936209446323930933335214619650277588138568504103427609283146728470455041360837549685364869161566863222680599110109210456299588352028330169041000166382937545505655464884266691630625402297821494221827392164049587946530563778771030124675514431,
		0.0479993885964583077281261798713460699543167134714936209446323930933335214619650277588138568504103427609283146728470455041360837549685364869161566863222680599110109210456299588352028330169041000166382937545505655464884266691630625402297821494221827392164049587946530563778771030124675514431,
		0.0475401657148303086622822069442231716408252512625387521584740318784735191312349586041971325618543660076682369564304738487584849740943805934034367382833518752314207901993991333786062812015195073547884746598535775062676699885664167707011249029305697669004958515813436770491520105115843742005,
		0.0475401657148303086622822069442231716408252512625387521584740318784735191312349586041971325618543660076682369564304738487584849740943805934034367382833518752314207901993991333786062812015195073547884746598535775062676699885664167707011249029305697669004958515813436770491520105115843742005,
		0.0469681828162100173253262857545810751998975284738125649829240886861900500181800807437012381630302198876925642461830694029139318555787845567143614289552410495903601238284556145544858090965965782916339169651505119399637862876053945518410353459767034026687936026945199383607112976484520939933,
		0.0469681828162100173253262857545810751998975284738125649829240886861900500181800807437012381630302198876925642461830694029139318555787845567143614289552410495903601238284556145544858090965965782916339169651505119399637862876053945518410353459767034026687936026945199383607112976484520939933,
		0.0462847965813144172959532492322611849696503075324468007778340818364698861774606986244241539105685321088517142947579291476238551538798963436740600968513359005801910700069462154098456091711311098901749803777735222026075473081311483686560830539773763176758567914860207820170792365910140063798,
		0.0462847965813144172959532492322611849696503075324468007778340818364698861774606986244241539105685321088517142947579291476238551538798963436740600968513359005801910700069462154098456091711311098901749803777735222026075473081311483686560830539773763176758567914860207820170792365910140063798,
		0.0454916279274181444797709969712690588873234618023998968168834081606504637618082102750954507142497706775055424364453740562113890878382679420378787427100982909191308430750899201141096789461078632697297091763378573830284133736378128577579722120264252594541491899441765769262904055702701625378,
		0.0454916279274181444797709969712690588873234618023998968168834081606504637618082102750954507142497706775055424364453740562113890878382679420378787427100982909191308430750899201141096789461078632697297091763378573830284133736378128577579722120264252594541491899441765769262904055702701625378,
		0.0445905581637565630601347100309448432940237999912217256432193286861948363377761089569585678875932857237669096941854082976565514031401996407675401022860761183118504326746863327792604337217763335682212515058414863183914930810334329596384915832703655935958010948424747251920190851700662833367,
		0.0445905581637565630601347100309448432940237999912217256432193286861948363377761089569585678875932857237669096941854082976565514031401996407675401022860761183118504326746863327792604337217763335682212515058414863183914930810334329596384915832703655935958010948424747251920190851700662833367,
		0.0435837245293234533768278609737374809227888974971180150532193925502569499020021803936448815937567079991401855477391110804568848623412043870399620479222000249538880795788245633051476595555730388360811011823841525667998427392843673284072004068821750061964976796287623004834501604656318714989,
		0.0435837245293234533768278609737374809227888974971180150532193925502569499020021803936448815937567079991401855477391110804568848623412043870399620479222000249538880795788245633051476595555730388360811011823841525667998427392843673284072004068821750061964976796287623004834501604656318714989,
		0.0424735151236535890073397679088173661655466481806496697314607722055245433487169327182398988553670128358787507582463602377168227019625334754497484024668087975720049504975593281010888062806587161032924284354938115463233015024659299046001504100674918329532481611571863222497170398830691222425,
		0.0424735151236535890073397679088173661655466481806496697314607722055245433487169327182398988553670128358787507582463602377168227019625334754497484024668087975720049504975593281010888062806587161032924284354938115463233015024659299046001504100674918329532481611571863222497170398830691222425,
		0.0412625632426235286101562974736380477399306355305474105429034779122704951178045914463267035032832336161816547420067160277921114474557623647771372636679857599931025531633255548770293397336318597716427093310378312957479805159734598610664983115148350548735211568465338522875618805992499897174,
		0.0412625632426235286101562974736380477399306355305474105429034779122704951178045914463267035032832336161816547420067160277921114474557623647771372636679857599931025531633255548770293397336318597716427093310378312957479805159734598610664983115148350548735211568465338522875618805992499897174,
		0.0399537411327203413866569261283360739186769506703336301114037026981570543670430333260307390357287606111017588757685176701688554806178713759519003171090525332423003042251947304213502522332118258365256241174986409729902714098049024753746340158430732115642207673265332738358717839602955875715,
		0.0399537411327203413866569261283360739186769506703336301114037026981570543670430333260307390357287606111017588757685176701688554806178713759519003171090525332423003042251947304213502522332118258365256241174986409729902714098049024753746340158430732115642207673265332738358717839602955875715,
		0.0385501531786156291289624969468089910127871122017180319662378854088005271323682681394418540442928363090545214563022868422017877042243007014244875098498616146404178795110038170109976252865902624380463581094085479557660525450020049773872343621719025128277593787164021147974906095237533202082,
		0.0385501531786156291289624969468089910127871122017180319662378854088005271323682681394418540442928363090545214563022868422017877042243007014244875098498616146404178795110038170109976252865902624380463581094085479557660525450020049773872343621719025128277593787164021147974906095237533202082,
		0.0370551285402400460404151018095833750834649453056563021747536272028091562122966687178302646649066832960609370472485057031765338738734008482025086366647963664178752038995704175623165041724901843573087856883034472545386037691055680911138721623610172486110313241291773258491882452773847899443,
		0.0370551285402400460404151018095833750834649453056563021747536272028091562122966687178302646649066832960609370472485057031765338738734008482025086366647963664178752038995704175623165041724901843573087856883034472545386037691055680911138721623610172486110313241291773258491882452773847899443,
		0.0354722132568823838106931467152459479480946310024100946926514848199381113651392962399922996268087884509143420993419937430515415557908457195618550238075571721209638845910166697234073788332647695349442265578792857058786796417110738673392400570019770741873271724201517438135222598792344040215,
		0.0354722132568823838106931467152459479480946310024100946926514848199381113651392962399922996268087884509143420993419937430515415557908457195618550238075571721209638845910166697234073788332647695349442265578792857058786796417110738673392400570019770741873271724201517438135222598792344040215,
		0.0338051618371416093915654821107254310210499263140045346675500650400323727745785853730452808963944098691936344225349051741060036935288424090581463711756382878498537611980973238606529148664990420534952057130296232922368792280098852092993207644225150541876980292972087619863453425206929192216,
		0.0338051618371416093915654821107254310210499263140045346675500650400323727745785853730452808963944098691936344225349051741060036935288424090581463711756382878498537611980973238606529148664990420534952057130296232922368792280098852092993207644225150541876980292972087619863453425206929192216,
		0.032057928354851553585467504347898716966221573881398062250169407854535275399124366530227987935629046729162364779969274126431870966979526186907589490002269660893281421728773647001279141626157958271220102615163092206489916992120482595587916535390136003611498634162765724522022671474313619317,
		0.032057928354851553585467504347898716966221573881398062250169407854535275399124366530227987935629046729162364779969274126431870966979526186907589490002269660893281421728773647001279141626157958271220102615163092206489916992120482595587916535390136003611498634162765724522022671474313619317,
		0.030234657072402478867974059819548659158281397768481241636026542045969161851838118212761980885178641520596873511042783163461341979185470882574743804555268086640389062237383427702813367624714014426121485626242067362445894463989335423458464954799181190120473168677930333898873084606011285311,
		0.030234657072402478867974059819548659158281397768481241636026542045969161851838118212761980885178641520596873511042783163461341979185470882574743804555268086640389062237383427702813367624714014426121485626242067362445894463989335423458464954799181190120473168677930333898873084606011285311,
		0.0283396726142594832275113052002373519812075841257543359907955185084500175712880712901834579816476269393013386531176072296695948860841466158639973753393323262188023471133258509422081952937349849822864752636994881600343083839805990853930436233762729622213044478376753949590318846038229829528,
		0.0283396726142594832275113052002373519812075841257543359907955185084500175712880712901834579816476269393013386531176072296695948860841466158639973753393323262188023471133258509422081952937349849822864752636994881600343083839805990853930436233762729622213044478376753949590318846038229829528,
		0.0263774697150546586716917926252251856755993308422457184496156736853021592428967790284780487213653480867620409279447766944383920384284787790772384251090745670478105870527396429136326932261251511732466974897397268573168068852344129736214469830280087710575094607457344820944885011053938108899,
		0.0263774697150546586716917926252251856755993308422457184496156736853021592428967790284780487213653480867620409279447766944383920384284787790772384251090745670478105870527396429136326932261251511732466974897397268573168068852344129736214469830280087710575094607457344820944885011053938108899,
		0.0243527025687108733381775504090689876499784155133784119819985685535536787083770723737264828464464223276155821319330210193549896426801083040150047332857692873011433649334477370145389017577189505240415125600908800786897201425473757275187332157593198572919772969833130729981971352463730545469,
		0.0243527025687108733381775504090689876499784155133784119819985685535536787083770723737264828464464223276155821319330210193549896426801083040150047332857692873011433649334477370145389017577189505240415125600908800786897201425473757275187332157593198572919772969833130729981971352463730545469,
		0.0222701738083832541592983303841550024229592905997594051455205429744914460867081990116647982811451138592401156680063927909718825845915896692701716212710541472344073624315399429951255221519263275095347974129106415903376085208797420439500915674568159744176912567285070988940509294826076696882,
		0.0222701738083832541592983303841550024229592905997594051455205429744914460867081990116647982811451138592401156680063927909718825845915896692701716212710541472344073624315399429951255221519263275095347974129106415903376085208797420439500915674568159744176912567285070988940509294826076696882,
		0.0201348231535302093723403167285438970895266801007919519220072343276769828211923597982299498416998597995443052252531684909219367615574440281549241161294448697202959593344989612626641188010558013085389491205901106884167596038790695150496733123662891637942237462337673353651179115491957031948,
		0.0201348231535302093723403167285438970895266801007919519220072343276769828211923597982299498416998597995443052252531684909219367615574440281549241161294448697202959593344989612626641188010558013085389491205901106884167596038790695150496733123662891637942237462337673353651179115491957031948,
		0.0179517157756973430850453020011193688971673570364158572977184273569247295870620984743089140579199272107974903016785911970727080884655646148340637373001805876560334052431930062983734905886704331100259778249929425439377011315288821865303197904492848823994202996722656114004109123107733596987,
		0.0179517157756973430850453020011193688971673570364158572977184273569247295870620984743089140579199272107974903016785911970727080884655646148340637373001805876560334052431930062983734905886704331100259778249929425439377011315288821865303197904492848823994202996722656114004109123107733596987,
		0.0157260304760247193219659952975397944260290098431565121528943932284210502164124556525745628476326997189475680077625258949765335021586482683126547283634704087193102431454662772463321304938516661086261262080252305539171654570677889578063634007609097342035360186636479612243231917699790225637,
		0.0157260304760247193219659952975397944260290098431565121528943932284210502164124556525745628476326997189475680077625258949765335021586482683126547283634704087193102431454662772463321304938516661086261262080252305539171654570677889578063634007609097342035360186636479612243231917699790225637,
		0.0134630478967186425980607666859556841084257719773496708184682785221983598894666268489697837056105038485845901773961664652581563686185523959473293683490869846700009741156668864960127745507806046701586435579547632680339906665338521813319281296935586498194608460412423723103161161922347608637,
		0.0134630478967186425980607666859556841084257719773496708184682785221983598894666268489697837056105038485845901773961664652581563686185523959473293683490869846700009741156668864960127745507806046701586435579547632680339906665338521813319281296935586498194608460412423723103161161922347608637,
		0.011168139460131128818590493019208135072778797816827287215251362273969701224836131369547661822970774719521543690039908073147476182135228738610704246958518755712518444434075738269866120460156365855324768445411463643114925829148750923090201475035559533993035986264487097245733097728698218563,
		0.011168139460131128818590493019208135072778797816827287215251362273969701224836131369547661822970774719521543690039908073147476182135228738610704246958518755712518444434075738269866120460156365855324768445411463643114925829148750923090201475035559533993035986264487097245733097728698218563,
		0.0088467598263639477230309146597306476951762660792204997984715769296110380005985367341694286322550520156167431790573509593010611842062630262878798782558974712042810219159674181580449655112696028911646066461502678711637780164986283350190669684468398617127841853445303466680698660632269500149,
		0.0088467598263639477230309146597306476951762660792204997984715769296110380005985367341694286322550520156167431790573509593010611842062630262878798782558974712042810219159674181580449655112696028911646066461502678711637780164986283350190669684468398617127841853445303466680698660632269500149,
		0.0065044579689783628561173603999812667711317610549523400952448792575685125717613068203530526491113296049409911387320826711045787146267036866881961532403342811327869183281273743976710008917886491097375367147212074243884772614562628844975421736416404173672075979097191581386023407454532945934,
		0.0065044579689783628561173603999812667711317610549523400952448792575685125717613068203530526491113296049409911387320826711045787146267036866881961532403342811327869183281273743976710008917886491097375367147212074243884772614562628844975421736416404173672075979097191581386023407454532945934,
		0.0041470332605624676352875357285514153133028192536848024628763661431834776690157393776820933106187137592011723199002845429836606307797425496666456172753165824787973801175029578301513761259541022471768825518482406145696380621686627285992715643614469568410535180218496973657001203470470418364,
		0.0041470332605624676352875357285514153133028192536848024628763661431834776690157393776820933106187137592011723199002845429836606307797425496666456172753165824787973801175029578301513761259541022471768825518482406145696380621686627285992715643614469568410535180218496973657001203470470418364,
		0.0017832807216964329472960791449719331799593472719279556695308063655858546954239803486698215802150348282744786016134857283616955449868451969230490863774274598030023211055562492709717566919237924255297982774711177411074145151155610163293142044147991553384925940046957893721166251082473659733,
		0.0017832807216964329472960791449719331799593472719279556695308063655858546954239803486698215802150348282744786016134857283616955449868451969230490863774274598030023211055562492709717566919237924255297982774711177411074145151155610163293142044147991553384925940046957893721166251082473659733,
	}
)
