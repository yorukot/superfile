---
title: Mengganti Komponen
description: Pelajari cara mengganti komponen bawaan Starlight untuk menambahkan elemen custom ke UI website dokumentasi Anda.
sidebar:
  badge: Baru
---

UI dan konfigurasi bawaan Starlight dirancang agar fleksibel dan dapat digunakan untuk berbagai macam jenis konten. Sebagian besar tampilan bawaan Starlight dapat disesuaikan dengan [CSS](/id/guides/css-and-tailwind/) dan [opsi konfigurasi](/id/guides/customization/).

Ketika Anda membutuhkan lebih dari apa yang telah disediakan, Starlight mendukung pembuatan komponen custom Anda sendiri untuk memperluas atau mengganti (sepenuhnya) komponen bawaannya.

## Kapan harus mengganti komponen

Mengganti komponen bawaan Starlight dapat berguna ketika:

- Anda ingin mengubah sebagian dari UI Starlight dengan cara yang tidak mungkin dilakukan dengan [custom CSS](/id/guides/css-and-tailwind/).
- Anda ingin mengubah bagaimana sebagian UI Starlight bekerja.
- Anda ingin menambahkan UI tambahan disamping UI Starlight yang sudah ada.

## Cara mengganti komponen

1. Pilih komponen Starlight yang ingin Anda ganti.
   Anda dapat menemukan daftar lengkap komponen di [Referensi Penggantian](/id/reference/overrides/).

   Contoh ini akan mengganti komponen [`SocialIcons`](/id/reference/overrides/#socialicons) Starlight pada bilah navigasi halaman.

2. Buat komponen Astro untuk menggantikan komponen Starlight tersebut.
   Contoh ini membuat tautan kontak.

   ```astro
   ---
   // src/components/EmailLink.astro
   import type { Props } from '@astrojs/starlight/props';
   ---

   <a href="mailto:houston@example.com">Kirim email ke Saya</a>
   ```

3. Beri tahu Starlight untuk menggunakan komponen custom Anda dalam opsi konfigurasi [`components`](/id/reference/configuration/#components) di `astro.config.mjs`:

   ```js {9-12}
   // astro.config.mjs
   import { defineConfig } from 'astro/config';
   import starlight from '@astrojs/starlight';

   export default defineConfig({
     integrations: [
       starlight({
         title: 'Website Dokumentasi Saya dengan Penggantian Komponen',
         components: {
           // Mengganti komponen `SocialIcons` bawaan.
           SocialIcons: './src/components/EmailLink.astro',
         },
       }),
     ],
   });
   ```

## Menggunakan ulang komponen bawaan

Anda dapat menggunakan komponen UI bawaan Starlight seperti yang Anda lakukan dengan komponen custom Anda sendiri: mengimpor dan merendernya dalam komponen custom Anda sendiri. Hal ini memungkinkan Anda mempertahankan semua elemen UI dasar Starlight dalam desain Anda, sambil menambahkan elemen UI tambahan bersama mereka.

Contoh di bawah ini menunjukkan sebuah komponen custom yang merender tautan e-mail bersama dengan komponen `SocialIcons` bawaan:

```astro {4,8}
---
// src/components/EmailLink.astro
import type { Props } from '@astrojs/starlight/props';
import Default from '@astrojs/starlight/components/SocialIcons.astro';
---

<a href="mailto:houston@example.com">Kirim email ke Saya</a>
<Default {...Astro.props}><slot /></Default>
```

Saat merender komponen bawaan dalam komponen custom:

- _Spread_ `Astro.props` ke dalamnya. Hal ini memastikan bahwa komponen tersebut menerima semua data yang diperlukan untuk merendernya.
- Tambahkan [`<slot />`](https://docs.astro.build/en/core-concepts/astro-components/#slots) di dalam komponen bawaan tersebut. Hal ini memastikan bahwa jika komponen tersebut menerima _child elements_, Astro tahu di mana merendernya.

## Menggunakan data halaman

Ketika mengganti komponen Starlight, implementasi custom Anda menerima objek standar `Astro.props` yang berisi semua data untuk halaman saat ini. Ini memungkinkan Anda menggunakan nilai-nilai ini untuk mengontrol bagaimana template komponen Anda merender.

Sebagai contoh, Anda dapat membaca nilai-nilai frontmatter halaman sebagai `Astro.props.entry.data`. Pada contoh berikut, komponen pengganti [`PageTitle`](/id/reference/overrides/#pagetitle) menggunakannya untuk menampilkan judul halaman saat ini:

```astro {5} "{title}"
---
// src/components/Title.astro
import type { Props } from '@astrojs/starlight/props';

const { title } = Astro.props.entry.data;
---

<h1 id="_top">{title}</h1>

<style>
  h1 {
    font-family: 'Comic Sans';
  }
</style>
```

Pelajari lebih lanjut tentang semua prop yang tersedia di [Referensi Penggantian](/id/reference/overrides/#component-props).

### Mengganti komponen hanya pada halaman tertentu

Penggantian komponen berlaku untuk semua halaman. Namun, Anda dapat merender secara kondisional menggunakan nilai dari `Astro.props` untuk menentukan kapan menampilkan UI custom Anda, kapan menampilkan UI bawaan Starlight, atau bahkan kapan menampilkan sesuatu yang benar-benar berbeda.

Pada contoh berikut, sebuah komponen yang menggantikan [`Footer`](/id/reference/overrides/#footer-1) Starlight menampilkan "Dibangun dengan Starlight 🌟" hanya di halaman beranda, dan menampilkan footer bawaan pada semua halaman lainnya:

```astro
---
// src/components/ConditionalFooter.astro
import type { Props } from '@astrojs/starlight/props';
import Default from '@astrojs/starlight/components/Footer.astro';

const isHomepage = Astro.props.slug === '';
---

{
  isHomepage ? (
    <footer>Dibangun dengan Starlight 🌟</footer>
  ) : (
    <Default {...Astro.props}>
      <slot />
    </Default>
  )
}
```

Pelajari lebih lanjut tentang merender secara kondisional di [Panduan Template Syntax Astro](https://docs.astro.build/en/core-concepts/astro-syntax/#dynamic-html).
