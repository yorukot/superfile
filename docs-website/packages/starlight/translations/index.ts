import { builtinI18nSchema } from '../schemas/i18n';
import cs from './cs.json';
import en from './en.json';
import es from './es.json';
import de from './de.json';
import ja from './ja.json';
import pt from './pt.json';
import fa from './fa.json';
import fr from './fr.json';
import gl from './gl.json';
import he from './he.json';
import id from './id.json';
import it from './it.json';
import nl from './nl.json';
import da from './da.json';
import tr from './tr.json';
import ar from './ar.json';
import nb from './nb.json';
import zh from './zh-CN.json';
import ko from './ko.json';
import sv from './sv.json';
import ro from './ro.json';
import ru from './ru.json';
import vi from './vi.json';
import uk from './uk.json';
import hi from './hi.json';
import zhTW from './zh-TW.json';
import pl from './pl.json';

const { parse } = builtinI18nSchema();

export default Object.fromEntries(
	Object.entries({
		cs,
		en,
		es,
		de,
		ja,
		pt,
		fa,
		fr,
		gl,
		he,
		id,
		it,
		nl,
		da,
		tr,
		ar,
		nb,
		zh,
		ko,
		sv,
		ro,
		ru,
		vi,
		uk,
		hi,
		'zh-TW': zhTW,
		pl,
	}).map(([key, dict]) => [key, parse(dict)])
);
