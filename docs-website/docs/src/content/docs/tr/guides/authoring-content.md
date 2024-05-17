---
title: Markdown'da İçerik Yazmak
description: Starlight'ın desteklediği Markdown sözdizimine genel bakış.
---

Starlight, `.md` uzantılı dosyalarda [Markdown](https://daringfireball.net/projects/markdown/) sözdizimini, [YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f) önbölümünde olduğu gibi başlık ve açıklama gibi metaverileri tanımlamak için destekler.

Markdown desteği ve kullanımı farklılık gösterebileceğinden dolayı, bu dosya formatlarını kullanıyorsanız [MDX dokümantasyonu](https://mdxjs.com/docs/what-is-mdx/#markdown) ya da [Markdoc dokümantasyonu](https://markdoc.dev/docs/syntax)'nu incelediğinizden emin olun.

## Satır İçi Stiller

Metin **kalın**, _italik_ ya da ~~üstü çizili~~ olabilir.

```md
Metin **kalın**, _italik_ ya da ~~üstü çizili~~ olabilir.
```

Başka bir sayfaya [bağlantı ekleyebilirsiniz](/tr/getting-started/).

```md
Başka bir sayfaya [bağlantı ekleyebilirsiniz](/tr/getting-started/).
```

Kesme işaretleri ile `satır için kodu` vurgulayabilirsiniz.

```md
Kesme işaretleri ile `satır için kodu` vurgulayabilirsiniz.
```

## Görseller

Starlight'ta görseller [Astro'nun kurulu optimize edilen dosya desteği](https://docs.astro.build/en/guides/assets/) ile kullanılır.

Markdown ve MDX, ekran okuyucular ve yardımcı teknolojilere yönelik alternatif metin içeren görselleri göstermek için Markdown sözdizimini destekler.

!["astro" metni içeren gezegen ve yıldızlar görseli](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
!["astro" metni içeren gezegen ve yıldızlar görseli](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)
```

Ayrıca, yerel olarak projenizde barındırılan görseller için ilişkili görsel dizin yolları desteklenir.

```md
// src/content/docs/page-1.md

![Uzayda bir roket](../../assets/images/rocket.svg)
```

## Başlıklar

Başlık kullanarak içerik yapınızı kurabilirsiniz. Markdown'daki başlıklar `#` sayısı ile satır başında oluşturulabilir.

### Starlight'ta sayfa içeriği yapısı nasıl kurulur

Starlight, sayfa başlığınızı en üst seviye başlık olarak kullanılacak şekilde yapılandırılmıştır ve içerik tablosunda "Genel Bakış" olarak yer alacaktır. Her sayfanın bir paragraf metniyle ve sayfa üstü başlığının `<h2>` ve alt seviyelerini kullanarak oluşturulmasını öneriyoruz:

```md
---
title: Markdown Rehberi
description: Starlight'ta Markdown nasıl kullanılır
---

Bu sayfa Starlight'ta Markdown'un nasıl kullanıldığını açıklar.

## Satır İçi Stiller

## Başlıklar
```

### Otomatik Başlık Bağlantıları

Markdown'da başlık kullanmak, otomatik olarak başlıklar için bağlantı oluşturur. Böylece sayfanızın belli bölümlerini direkt olarak bağlantılandırabilirsiniz:

```md
---
title: Sayfa İçeriğim
description: Starlight'ın kurulu bağlantıları nasıl kullanılır
---

## Giriş

[Görüşümü](#görüş) aynı sayfanın aşağısına iliştirebilirim.

## Görüş

`https://my-site.com/page1/#introduction` giriş bölümüme direkt olarak yönlendirir.
```

Seviye 2 (`<h2>`) ve Seviye 3 (`<h3>`) başlıklar otomatik olarak içerik tablosunda görünecektir.

## Ara Bölümler

Ara bölümler, sayfanın ana içeriğinin yanında ikincil bilgi gösterimi için kullanışlıdır.

Starlight ara bölümleri oluşturmak için özel Markdown sözdizimi sunar. Ara bölüm blokları üç adet iki nokta üst üste'nin `:::` içeriği sarmalamasıyla kullanılır ve tip olarak `note`,`tip`, `caution` ya da `danger` kullanılabilir.

Herhangi bir Markdown içerik tipini ara bölümü içerisine yerleştirebilirsiniz, ancak ara bölümler kısa ve öz içerikler için biçilmiş kaftandır.

### Note Ara bölümü

:::note
Starlight, [Astro](https://astro.build/) ile oluşturulmuş bir dokümantason website oluşturma aracıdır. Bu komutla başlayabilirsiniz:

```sh
npm create astro@latest -- --template starlight
```

:::

````md
:::note
Starlight, [Astro](https://astro.build/) ile oluşturulmuş bir dokümantason website oluşturma aracıdır. Bu komutla başlayabilirsiniz:

```sh
npm create astro@latest -- --template starlight
```

:::
````

### Özel ara bölümler

Ara bölümler için ara bölüm tipinin tanımından hemen sonra, köşeli parantez arasında olacak şekilde ara bölümlerinizi özelleştirebilirsiniz (örn.`:::tip[Bunu biliyor musun?]`).

:::tip[Bunu biliyor musun?]
Astro [Ada Mimarisi”](https://docs.astro.build/en/concepts/islands/) ile daha hızlı websitesi oluşturmana yardımcı olur.
:::

```md
:::tip[Bunu biliyor musun?]
Astro [Ada Mimarisi”](https://docs.astro.build/en/concepts/islands/) ile daha hızlı websitesi oluşturmana yardımcı olur.
:::
```

### Diğer Ara bölümler

Uyarı ve tehlike ara bölümleri, kullanıcıların dikkatini gözden kaçabilecek detaylara çekmek için kullanışlıdır.
Bunları çok kullandığınızı farkederseniz, dokümanınızın yeniden oluşturulmasına gerek kalmayacağının işareti olabilir.

:::caution[Uyarı]
Harika bir dokümantasyon sitesi istediğine emin değilsen, [Starlight](/tr/) kullanmadan önce iki kez düşün.
:::

:::danger[Tehlike]
Yardımcı Starlight özellikleri sayesinde kullanıcılarınız daha kolay ürün bulabilir ve daha üretken olabilir.

- Yönlendirmeyi temizle
- Kullanıcı-yapılandırmalı renk teması
- [i18n desteği](/tr/guides/i18n/)

:::

```md
:::caution
Harika bir dokümantasyon sitesi istediğine emin değilsen, [Starlight](/tr/) kullanmadan önce iki kez düşün.
:::

:::danger
Yardımcı Starlight özellikleri sayesinde kullanıcılarınız daha kolay ürün bulabilir ve daha üretken olabilir.

- Yönlendirmeyi temizle
- Kullanıcı-yapılandırmalı renk teması
- [i18n desteği](/tr/guides/i18n/)

:::
```

## Blok Alıntılar

> Bu, genelde başka bir belge ya da kişiden alıntılanan bir blok alıntıdır.
>
> Blok alıntılar her satırda `>` ile başlar.

```md
> Bu, genelde başka bir belge ya da kişiden alıntılanan bir blok alıntıdır.
>
> Blok alıntılar her satırda `>` ile başlar.
```

## Kod Blokları

Kod bloğu, başında ve sonunda üç kesme işaretinin arasında kalan <code>```</code> bir bloktur. Üç kesme işaretiyle başladıktan hemen sonra göstermek istediğiniz programlama dilini belirtebilirsiniz.

```js
// Sözdizimi vurgulamalı Javascript kodu.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```

````md
```js
// Sözdizimi vurgulamalı Javascript kodu.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```
````

```md
Uzun, tek satırlı kod bloğu alt satıra geçmemelidir. Çok uzunsa yatay kaydırma olmalıdır. Bu satır, yatay kaydırma çubuğunun görünmesi için yeterince uzun olmalıdır.
```

## Diğer ortak Markdown Özellikleri

Starlight, liste ve tablo gibi diğer tüm Markdown yazım sözdizimini destekler. [Markdown Rehberi'nden Markdown Kopya Kağıdı](https://www.markdownguide.org/cheat-sheet/)'na tüm Markdown sözdizimi elemanlarına hızlı bir genel bakış için göz atın.
