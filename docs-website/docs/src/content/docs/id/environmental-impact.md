---
title: Dokumentasi yang ramah lingkungan
description: Pelajari bagaimana Starlight bisa membantu Anda membangun website dokumentasi yang lebih ramah lingkungan dan mengurangi jejak karbon Anda.
---

Estimasi dampak dari industri web terhadap iklim berkisar antara [2%][sf] hingga [4% dari emisi karbon global ][bbc], kira-kira setara dengan emisi industri penerbangan. Ada banyak faktor kompleks dalam menghitung dampak ekologis sebuah website, namun panduan ini mencakup beberapa tips untuk mengurangi jejak lingkungan dari website dokumentasi Anda.

Berita baiknya adalah, memilih Starlight adalah awal yang baik. Menurut Website Carbon Calculator, website ini [lebih bersih daripada 99% website-website yang telah diuji][sl-carbon], menghasilkan 0,01g CO₂ per kunjungan halaman.

## Berat halaman

Semakin banyak data yang ditransfer oleh sebuah halaman web, semakin banyak sumber daya energi yang diperlukan.
Pada bulan April 2023, nilai median dari banyaknya data yang harus di-_download_ user ketika mengakses sebuah halaman website adalah lebih dari 2.000 KB berdasarkan [data dari HTTP Archive][http].

Starlight membangun halaman-halaman yang seringan mungkin. Sebagai contoh, pada kunjungan pertama, pengguna hanya perlu mengunduh kurang dari 50 KB data yang telah dikompresi — hanya 2,5% dari nilai median HTTP Archive. Dengan strategi caching yang baik, kunjungan selanjutnya dapat mengunduh hanya sekitar 10 KB.

### Gambar

Meskipun Starlight memberikan basis yang baik, gambar yang Anda tambahkan ke halaman dokumentasi Anda dapat dengan cepat meningkatkan berat halaman Anda.
Starlight menggunakan [dukungan aset yang dioptimalkan][assets] dari Astro untuk mengoptimalkan gambar lokal dalam file Markdown dan MDX Anda.

### Komponen UI

Komponen yang dibangun dengan _UI frameworks_ seperti React atau Vue dapat dengan mudah menambahkan banyak JavaScript ke halaman.
Karena Starlight dibangun di atas Astro, komponen seperti ini secara default tidak memuat JavaScript di sisi klien berkat [Astro Islands][islands].

### Caching

_Caching_ digunakan untuk mengontrol berapa lama browser menyimpan dan menggunakan kembali data yang telah diunduh sebelumnya.
Strategi caching yang baik memastikan bahwa pengguna mendapatkan konten baru sesegera mungkin ketika ada perubahan, tetapi juga menghindari pengunduhan yang tidak perlu dari konten yang sama berulang kali ketika konten tersebut tidak mengalami perubahan.

Cara paling umum untuk mengonfigurasi caching adalah dengan menggunakan [`Cache-Control` HTTP header][cache].
Ketika menggunakan Starlight, Anda dapat mengatur waktu _cache_ yang lama untuk semua yang ada di direktori /\_astro/.
Direktori ini berisi CSS, JavaScript, dan aset lainnya yang dapat di-cache secara permanen, mengurangi pengunduhan yang tidak perlu:

```
Cache-Control: public, max-age=604800, immutable
```

Cara mengkonfigurasi caching tergantung pada penyedia hosting website Anda. Misalnya, Vercel menerapkan strategi caching ini untuk Anda tanpa ada konfigurasi yang diperlukan, sementara Anda dapat mengatur [header kustom untuk Netlify][ntl-headers] dengan menambahkan file `public/_headers` ke proyek Anda:

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## Konsumsi daya

Cara sebuah halaman web dibangun dapat mempengaruhi besarnya daya yang dibutuhkan untuk menjalankannya di perangkat pengguna.
Dengan menggunakan JavaScript yang minimal, Starlight mengurangi jumlah daya pemrosesan yang dibutuhkan oleh telepon, tablet, atau komputer pengguna untuk memuat dan merender halaman.

Perhatikan saat menambahkan fitur seperti skrip pelacakan analitik atau konten yang kaya akan JavaScript seperti video yang disematkan, karena hal ini dapat meningkatkan penggunaan daya halaman.
Jika Anda memerlukan analitik, pertimbangkan untuk memilih opsi yang lebih ringan seperti [Cabin][cabin], [Fathom][fathom], atau [Plausible][plausible].
Penyisipan video seperti YouTube dan Vimeo dapat ditingkatkan dengan menunggu [pemuatan video saat ada interaksi pengguna][lazy-video].
_Package_ seperti [astro-embed][embed] dapat membantu untuk layanan umum.

:::tip[Tahukah Anda?]
_Parsing_ dan kompilasi JavaScript adalah salah satu tugas yang paling mahal bagi browser.
Dibandingkan dengan merender gambar JPEG dengan ukuran yang sama, [pemrosesan JavaScript dapat memakan waktu lebih dari 30 kali lebih lama][cost-of-js].
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## Hosting

Dimana website di-_hosting_ dapat memiliki dampak besar terhadap seberapa ramah lingkungan website dokumentasi Anda.
Pusat data dan rumah server dapat memiliki dampak ekologis yang besar, termasuk konsumsi listrik yang tinggi dan penggunaan air yang intensif.

Memilih penyedia hosting yang menggunakan energi terbarukan berarti emisi karbon yang lebih rendah untuk website Anda. [Green Web Directory][gwb] adalah salah satu alat yang dapat membantu Anda menemukan perusahaan hosting yang ramah lingkungan.

[gwb]: https://www.thegreenwebfoundation.org/directory/

## Perbandingan

Tertarik bagaimana perbandingannya dengan _framework_ dokumentasi lainnya? Tes ini dengan [Website Carbon Calculator][wcc] membandingkan halaman-halaman serupa yang dibangun dengan _tool_ yang berbeda.

| Framework                   | CO₂ per kunjungan halaman |
| --------------------------- | ------------------------- |
| [Starlight][sl-carbon]      | 0.01g                     |
| [VitePress][vp-carbon]      | 0.05g                     |
| [Docus][dc-carbon]          | 0.05g                     |
| [Sphinx][sx-carbon]         | 0.07g                     |
| [MkDocs][mk-carbon]         | 0.10g                     |
| [Nextra][nx-carbon]         | 0.11g                     |
| [docsify][dy-carbon]        | 0.11g                     |
| [Docusaurus][ds-carbon]     | 0.24g                     |
| [Read the Docs][rtd-carbon] | 0.24g                     |
| [GitBook][gb-carbon]        | 0.71g                     |

<small>Data dikumpulkan pada 14 Mei 2023. Klik link untuk melihat angka terkini.</small>

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

## Sumber Tambahan

### Tools

- [Website Carbon Calculator][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [WebPageTest Carbon Control](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### Articles and presentasi

- [“Building a greener web”](https://youtu.be/EfPoOt7T5lg), talk by Michelle Barker
- [“Sustainable Web Development Strategies Within An Organization”](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/), article by Michelle Barker
- [“A sustainable web for everyone”](https://2021.stateofthebrowser.com/speakers/tom-greenwood/), talk by Tom Greenwood
- [“How Web Content Can Affect Power Usage”](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/), article by Benjamin Poulain and Simon Fraser

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/en/guides/assets/
[islands]: https://docs.astro.build/en/concepts/islands/
[wcc]: https://www.websitecarbon.com/
