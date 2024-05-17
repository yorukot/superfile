---
title: Menulis Konten dalam Format Markdown
description: Gambaran tentang sintaks Markdown yang didukung oleh Starlight.
---

Starlight mendukung seluruh variasi sintaks [Markdown](https://daringfireball.net/projects/markdown/) dalam file `.md` serta menggunakan frontmatter [YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f) untuk mendefinisikan metadata seperti judul dan deskripsi.

Pastikan untuk mengecek [dokumentasi MDX](https://mdxjs.com/docs/what-is-mdx/#markdown) atau [dokumentasi Markdoc](https://markdoc.dev/docs/syntax) jika menggunakan format file tersebut, karena dukungan dan penggunaan Markdown dapat berbeda.

## Inline styles

Teks bisa **bold**, _italic_, atau ~~strikethrough~~.

```md
Teks bisa **bold**, _italic_, atau ~~strikethrough~~.
```

Anda dapat [menautkan ke halaman lain](/id/getting-started/).

```md
Anda dapat [menautkan ke halaman lain](/id/getting-started/).
```

Anda dapat menandakan `inline code` dengan _backticks_.

```md
Anda dapat menandakan `inline code` dengan _backticks_.
```

## Gambar

Gambar dalam Starlight menggunakan [dukungan aset teroptimalkan bawaan Astro](https://docs.astro.build/en/guides/assets/).

Markdown dan MDX mendukung sintaks Markdown untuk menampilkan gambar yang mencakup teks alternatif untuk pembaca layar dan teknologi assistif.

![Ilustrasi planet dan bintang dengan kata “astro“](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
![Ilustrasi planet dan bintang dengan kata “astro“](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)
```

_Relative paths_ juga didukung untuk gambar yang disimpan secara lokal di proyek anda.

```md
// src/content/docs/page-1.md

![Roket di luar angkasa](../../assets/images/rocket.svg)
```

## Judul

Anda dapat menyusun konten dengan menggunakan judul. Judul dalam Markdown ditandai dengan sejumlah `#` di awal baris.

### Bagaimana cara menyusun konten halaman di Starlight

Starlight dikonfigurasi untuk secara otomatis menggunakan judul halaman Anda sebagai judul tingkat atas dan akan menyertakan judul "Ringkasan" di bagian atas daftar isi setiap halaman. Kami merekomendasikan memulai setiap halaman dengan konten paragraf biasa dan menggunakan judul di dalam halaman dari `<h2>` ke bawah:

```md
---
title: Panduan Markdown
description: Cara menggunakan Markdown di Starlight
---

Halaman ini menjelaskan cara menggunakan Markdown di Starlight.

## Inline Styles

## Judul
```

### Automatic heading anchor links

Menggunakan judul dalam Markdown secara otomatis akan memberi Anda _anchor links_ sehingga Anda dapat langsung menautkan ke bagian-bagian tertentu dari halaman Anda:

```md
---
title: Halaman Konten Saya
description: Cara menggunakan _anchor links_ bawaan Starlight
---

## Pengantar

Saya dapat menautkan ke [kesimpulan saya](#kesimpulan) di bagian bawah halaman yang sama.

## Kesimpulan

`https://situs-saya.com/halaman1/#pengantar` langsung menuju ke Pengantar saya.
```

Judul Level 2 (`<h2>`) dan Level 3 (`<h3>`) akan secara otomatis muncul di daftar isi halaman.

## Asides

_Asides_ (juga sering disebut sebagai _“admonitions”_ atau “_callouts”_) berguna untuk menampilkan informasi sekunder di samping konten utama halaman.

Starlight menyediakan sintaks Markdown kustom untuk merender _asides_. Blok _asides_ ditandai dengan sepasang tiga titik dua `:::` untuk melingkupi konten Anda, dan dapat berjenis `note`, `tip`, `caution`, atau `danger`.

Anda dapat menyusun berbagai jenis konten Markdown lainnya di dalam sebuah _asides_, tetapi _asides_ lebih cocok untuk potongan konten yang pendek dan padat.

### Catatan Asides

:::note
Starlight adalah toolkit website dokumentasi yang dibangun dengan [Astro](https://astro.build/). Anda dapat memulai dengan perintah ini:

```sh
npm create astro@latest -- --template starlight
```

:::

````md
:::note
Starlight adalah toolkit website dokumentasi yang dibangun dengan [Astro](https://astro.build/). Anda dapat memulai dengan perintah ini:

```sh
npm create astro@latest -- --template starlight
```

:::
````

### Judul Asides Kustom

Anda dapat menentukan judul kustom untuk _asides_ dengan menambahkan tanda kurung siku setelah jenis _asides-nya_, misalnya `:::tip[Apakah Anda tahu?]`.

:::tip[Apakah Anda tahu?]
Astro membantu Anda membangun website lebih cepat dengan [“Islands Architecture”](https://docs.astro.build/en/concepts/islands/).
:::

```md
:::tip[Apakah Anda tahu?]
Astro membantu Anda membangun website lebih cepat dengan [“Islands Architecture”](https://docs.astro.build/en/concepts/islands/).
:::
```

### Jenis Asides Lainnya

_Asides_ berjenis _caution_ dan _danger_ berguna untuk menarik perhatian pengguna pada detail-detail yang mungkin membuat mereka bingung.
Jika Anda sering menggunakan ini, mungkin juga pertanda bahwa hal yang Anda dokumentasikan sepertinya bisa di-desain ulang.

:::caution
Jika Anda tidak yakin ingin membuat situs dokumentasi yang menakjubkan, pikirkan dua kali sebelum menggunakan [Starlight](/id/).
:::

:::danger
Pengguna Anda mungkin lebih produktif dan menemukan produk Anda lebih mudah digunakan berkat fitur-fitur Starlight yang membantu.

- Navigasi yang jelas
- Tema warna yang dapat dikonfigurasi oleh pengguna
- [Dukungan i18n](/id/guides/i18n/)

:::

```md
:::caution
Jika Anda tidak yakin ingin membuat situs dokumen yang menakjubkan, pikirkan dua kali sebelum menggunakan [Starlight](/id/).
:::

:::danger
Pengguna Anda mungkin lebih produktif dan menemukan produk Anda lebih mudah digunakan berkat fitur-fitur Starlight yang membantu.

- Navigasi yang jelas
- Tema warna yang dapat dikonfigurasi oleh pengguna
- [Dukungan i18n](/id/guides/i18n/)

:::
```

## Blockquote

> Ini adalah blockquote, yang biasanya digunakan saat mengutip orang lain atau dokumen lain.
>
> Blockquotes ditandai dengan tanda `>` di awal setiap barisnya.

```md
> Ini adalah blockquote, yang biasanya digunakan saat mengutip orang lain atau dokumen lain.
>
> Blockquotes ditandai dengan tanda `>` di awal setiap barisnya.
```

## Code blocks

_Code blocks_ ditandai dengan blok tiga tanda kutip terbalik <code>```</code> di awal dan akhir. Anda dapat menunjukkan bahasa pemrograman yang digunakan setelah tanda kutip terbalik pembuka.

```js
// Kode javascript dengan syntax highlighting.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```

````md
```js
// Kode javascript dengan syntax highlighting.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```
````

```md
Kode satu baris tunggal yang panjang sebaiknya tidak di-wrap. Kode tersebut harus menggulir secara horizontal jika terlalu panjang. Baris ini sudah cukup panjang untuk mencontohkan hal yang dimaksud.
```

## Fitur Umum Markdown Lainnya

Starlight mendukung penulisan semua sintaks Markdown lainnya, seperti daftar dan tabel. Lihat [Markdown Cheat Sheet dari The Markdown Guide](https://www.markdownguide.org/cheat-sheet/) untuk penjelasan singkat tentang semua sintaks elemen Markdown.

## Konfigurasi Markdown dan MDX Lanjutan

Starlight menggunakan Markdown dan renderer MDX Astro yang dibangun berdasarkan remark dan rehype. Anda dapat menambahkan dukungan untuk sintaks dan perilaku khusus dengan menambahkan `remarkPlugins` atau `rehypePlugins` di file konfigurasi Astro Anda. Lihat [“Configuring Markdown and MDX”](https://docs.astro.build/en/guides/markdown-content/#configuring-markdown-and-mdx) dalam dokumentasi Astro untuk mempelajari lebih lanjut.
