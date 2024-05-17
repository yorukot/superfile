---
title: Doğa Dostu Dokümanlar
description: Starlight'ın daha çevreci dokümantasyon sitesi oluşturmanda nasıl yardımcı olacağını ve karbon ayakizini nasıl düşüreceğini öğren.
---

Web endüstrisinin iklime etkisinin [%2][sf] ila [4% arasında küresel karbon emisyonu][bbc]'na sahip olduğu tahmin ediliyor; kabaca havayolları endüstrisindeki emisyon oranında.
Bir websitesinin ekolojik etkisini hesaplamada çok fazla karmaşık etken var, ancak bu rehber dokümantasyon sitenizin çevresel ayakizini düşürmek için birkaç ipucu içerir.

İyi haber şu ki, Starlight'ı seçmek mükemmel bir başlangıç.
Bu site, Carbon Calculator web sitesine göre sayfa başına ziyarette 0.01 gr CO₂ üretimiyle [test edilmiş web sayfaları arasında %99 daha çevreci][sl-carbon].

## Sayfa Büyüklüğü

Bir web sayfası ne kadar veri transfer ederse, ihtiyaç duyacağı enerji o kadar yüksek olur.
Nisan 2023'te [HTTP Archive'daki veri][http]'ye göre orta büyüklükte web sayfası, kullanıcıdan 2,000 KB'tan daha fazla veri indirmeye ihtiyaç duydu.

Starlight, sayfaları mümkün olan en küçük boyutta oluşturur.
Örneğin, ilk ziyarette kullanıcı 50 KB'tan daha az sıkıştırılmış veriyi indirecektir - sadece HTTP Archive medyan değerinin %2.5'i kadar.
İyi bir ön belleğe alma stratejisi ile takip eden gezinmelerde 10 KB kadar küçük veri indirir.

### Görseller

Starlight iyi bir temel sunmasına rağmen, dokümantasyon sayfalarına eklediğiniz görseller sayfa büyüklüğünü hızlıca artırabilir.
Starlight, Astro'nun [optimize edilmiş varlık desteği][assets] yardımıyla Markdown ve MDX dosyalarındaki yerel görselleri optimize etmek için kullanır.

### Arayüz Bileşenleri

React ya da Vue gibi arayüz kütüphaneleri ile oluşturulmuş bileşenler, sayfaya büyük boyutta Javascript ekleyebilir.
[Astro Adaları][islands] ile oluşturulmuş bileşenlerin **varsayılan olarak sıfır tarayıcı-tarafı Javascript** yüklemesi nedeniyle Starlight, Astro üzerine kurulmuştur.

### Ön Bellekleme

Ön bellekleme, halihazırda yüklenmiş verilerin tekrar kullanımı ve tarayıcının ne kadar süre bu veriyi tutacağını kontrol etmek için kullanılır.
İyi bir ön bellekleme stratejisi; kullanıcının yeni içeriği değiştikten sonra mümkün olan en kısa sürede almasını sağlar, fakat buna ek olarak içerik değişmediğinde tekrar tekrar aynı içeriği yüklemesini önler.

[`Cache-Control` HTTP header][cache], ön belleklemeyi yapılandırmanın sık kullanılan bir yoludur.
Starlight kullanarak, `/_astro/` dizini içindeki her şey için uzun ön bellekleme süresi ayarlayabilirsiniz.
Bu dizin CSS, Javascript ve diğer paketlenmiş dosyaları içerir ve güvenle sonsuza dek ön belleğe alınarak gereksiz indirmeleri azaltır:

```
Cache-Control: public, max-age=604800, immutable
```

Ön belleklemenin nasıl yapılandırılacağı web sunucunuza bağlıdır. Örneğin; Vercel önbellekleme stratejisini herhangi bir yapılandırmaya gerek kalmadan uygularken, Netlify'da `public/_headers` altına [Netlify için özel header'lar][ntl-headers] ile ayarlayabilirsiniz:

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## Güç tüketimi

Bir web sitesinin nasıl oluşturulduğu, web sitesinin kullanıcı cihazında ne kadar güç tüketeceğini etkiler.
En az Javascript kullanarak, Starlight kullanıcının telefon, tablet ve bilgisayarında yükleme ve sayfayı çizme sürecinde ihtiyaç duyulan enerjiyi düşürür.

Analitik takip scriptleri ya da gömülü video gibi ağır Javascript içeriklerini eklerken, bu özelliklerin sayfanızın güç tüketimini artıracağını aklınızda tutun.
Analitik takip script'ine ihtiyacınız varsa, [Cabin][cabin], [Fathom][fathom], ya da [Plausible][plausible] gibi daha hafif seçenekleri göz önünde bulundurun.
Youtube ve Vimeo gibi gömülü videolar [kullanıcı etkileşimi ile videoyu yükleme][lazy-video]. yöntemiyle geliştirilebilir.
[`astro-embed`][embed] gibi paketler ortak hizmetler için yardımcı olabilir.

:::tip[Bunu biliyor musun?]
Javascript'in ayrıştırılması ve derlenmesi tarayıcıların yapması gereken en maliyetli görevlerden biri.
Aynı boyutta JPEG görsel çizimi ile kıyaslandığında, [JavaScript 30 kez daha uzun işleme süresi alabilir][cost-of-js].
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## Barındırma

Web sayfasının nerede barındırıldığı, dokümantasyon sitenizin ne kadar çevre dostu olduğu konusunda büyük etkisi vardır.
Veri merkezleri ve sunucu çiftlikleri büyük ekolojik etkiye sahip olabilir, yüksek elektrik tüketimi ve büyük oranda su kullanımı dahil.

Yenilenebilir enerji kullanan barındırma opsiyonunu seçmek, sitenizin daha az karbon emisyonuna neden olacağı anlamına gelir. [Green Web Directory][gwb] size bu konuda barındırma şirketleri bulmanıza yardımcı olacak bir araçtır.

[gwb]: https://www.thegreenwebfoundation.org/directory/

## Karşılaştırmalar

Diğer dokümantasyon çerçeveleri ile nasıl kıyaslandığını merak ediyor musun?
[Website Carbon Calculator][wcc] ile yapılan testlerle, farklı araçlarla oluşturulmuş benzer sayfaları karşılaştırın.

| Çerçeve                     | sayfa başına ziyarette CO₂ |
| --------------------------- | -------------------------- |
| [Starlight][sl-carbon]      | 0.01gr                     |
| [VitePress][vp-carbon]      | 0.05gr                     |
| [Docus][dc-carbon]          | 0.05gr                     |
| [Sphinx][sx-carbon]         | 0.07gr                     |
| [MkDocs][mk-carbon]         | 0.10gr                     |
| [Nextra][nx-carbon]         | 0.11gr                     |
| [docsify][dy-carbon]        | 0.11gr                     |
| [Docusaurus][ds-carbon]     | 0.24gr                     |
| [Read the Docs][rtd-carbon] | 0.24gr                     |
| [GitBook][gb-carbon]        | 0.71gr                     |

<small>14 Mayıs 2023 tarihli veriler. Güncel bilgileri görmek için linke tıkla.</small>

[sl-carbon]: https://www.websitecarbon.com/website/starlight-astro-build-getting-started/
[vp-carbon]: https://www.websitecarbon.com/website/vitepress-dev-guide-what-is-vitepress/
[dc-carbon]: https://www.websitecarbon.com/website/docus-dev-introduction-getting-started/
[sx-carbon]: https://www.websitecarbon.com/website/sphinx-doc-org-en-master-usage-quickstart-html/
[mk-carbon]: https://www.websitecarbon.com/website/mkdocs-org-getting-started/
[nx-carbon]: https://www.websitecarbon.com/website/nextra-site-docs-docs-theme-start/
[dy-carbon]: https://www.websitecarbon.com/website/docsify-js-org/
[ds-carbon]: https://www.websitecarbon.com/website/docusaurus-io-docs/
[rtd-carbon]: https://www.websitecarbon.com/website/docs-readthedocs-io-en-stable-index-html/
[gb-carbon]: https://www.websitecarbon.com/website/docs-gitbook-com/

## Diğer Kaynaklar

### Araçlar

- [Website Carbon Calculator][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [WebPageTest Carbon Control](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### Makaleler ve Konuşmalar

- [“Building a greener web”](https://youtu.be/EfPoOt7T5lg), Michelle Barker'ın konuşması
- [“Sustainable Web Development Strategies Within An Organization”](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/), Michelle Barker'ın makalesi
- [“A sustainable web for everyone”](https://2021.stateofthebrowser.com/speakers/tom-greenwood/), Tom Greenwood'un konuşması
- [“How Web Content Can Affect Power Usage”](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/), Benjamin Poulain and Simon Fraser'ın makalesi

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/en/guides/assets/
[islands]: https://docs.astro.build/en/concepts/islands/
[wcc]: https://www.websitecarbon.com/
