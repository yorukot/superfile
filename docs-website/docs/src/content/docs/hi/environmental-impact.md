---
title: पर्यावरण-अनुकूल दस्तावेज़
description: जानें कि कैसे Starlight आपको हरित दस्तावेज़ीकरण साइट बनाने और आपके कार्बन निशान को कम करने में मदद कर सकता है।
---

वेब उद्योग के जलवायु प्रभाव का अनुमान वैश्विक कार्बन उत्सर्जन के [2%][sf] से [4% के बीच][bbc] है, जो लगभग एयरलाइन उद्योग के उत्सर्जन के बराबर है।
किसी वेबसाइट के पारिस्थितिक प्रभाव की गणना करने में कई जटिल कारक होते हैं, लेकिन इस मार्गदर्शिका में आपके दस्तावेज़ साइट के पर्यावरणीय पदचिह्न को कम करने के लिए कुछ युक्तियां शामिल हैं।

अच्छी खबर यह है कि Starlight को चुनना एक बेहतरीन शुरुआत है।
Website Carbon Calculator के अनुसार, यह साइट [परीक्षण किए गए 99% वेब पेजों की तुलना में अधिक स्वच्छ है][sl-carbon], जो प्रति पेज भेंट में 0.01 ग्राम CO₂ का उत्पादन करती है।

## पेज का वजन

एक वेब पेज जितना अधिक डेटा स्थानांतरित करता है, उसे उतने ही अधिक शक्ति संसाधनों की आवश्यकता होती है।
एप्रिल 2023 में, [HTTP आर्काइव के डेटा][http] के अनुसार, औसत वेब पेज के लिए उपयोगकर्ता को 2,000 KB से अधिक डाउनलोड करने की आवश्यकता थी

Starlight ऐसे पेज बनाता है जो यथासंभव हल्के होते हैं।
उदाहरण के लिए, पहले दर्शन पर, एक उपयोगकर्ता 50 KB से कम संपीड़ित डेटा डाउनलोड करेगा - HTTP संग्रह माध्य का केवल 2.5%।
एक अच्छी कैशिंग रणनीति के साथ, बाद के नेविगेशन कम से कम 10 KB तक डाउनलोड हो सकते हैं।

### छवियाँ

जबकि Starlight एक अच्छी आधार रेखा प्रदान करता है, आपके द्वारा अपने दस्तावेज़ पेजों में जोड़ी गई छवियां आपके पेज का वजन तेजी से बढ़ा सकती हैं।
Starlight आपके Markdown और MDX फ़ाइलों में स्थानीय छवियों को अनुकूलित करने के लिए Astro के [अनुकूलित संपत्ति समर्थन][assets] का उपयोग करता है।

### UI अवयव

React या Vue जैसे UI फ्रेमवर्क के साथ निर्मित अवयव आसानी से एक पेज पर बड़ी मात्रा में Javascript जोड़ सकते हैं।
क्योंकि Starlight Astro पर बनाया गया है, [Astro द्वीप समूह][islands] की बदौलत इन जैसे अवयव डिफ़ॉल्ट रूप से **शून्य client-side Jacascript** लोड करते हैं।

### कैशिंग

कैशिंग का उपयोग यह नियंत्रित करने के लिए किया जाता है कि कोई ब्राउज़र पहले से डाउनलोड किए गए डेटा को कितनी देर तक संग्रहीत कर सकता है और उसका पुन: उपयोग करता है।
एक अच्छी कैशिंग रणनीति यह सुनिश्चित करती है कि उपयोगकर्ता को नई सामग्री बदलने पर जल्द से जल्द मिल जाए, लेकिन जब सामग्री नहीं बदली है तो उसे बार-बार डाउनलोड करने से भी बचा जा सकता है।

कैशिंग को कॉन्फ़िगर करने का सबसे आम तरीका [`Cache-Control` HTTP header][cache] है।
Starlight का उपयोग करते समय, आप `/_astro/` निर्देशिका में हर चीज़ के लिए एक लंबा cache समय निर्धारित कर सकते हैं।
इस निर्देशिका में CSS, JavaScript और अन्य बंडल संपत्तियां शामिल हैं जिन्हें अनावश्यक डाउनलोड को कम करते हुए हमेशा के लिए सुरक्षित रूप से cache किया जा सकता है:

```
Cache-Control: public, max-age=604800, immutable
```

कैशिंग को कैसे कॉन्फ़िगर करें यह आपके वेब होस्ट पर निर्भर करता है। उदाहरण के लिए, Vercel बिना किसी कॉन्फ़िगरेशन की आवश्यकता के आपके लिए इस कैशिंग रणनीति को लागू करता है, जबकि आप अपने परियोजना में `public/_headers` फ़ाइल जोड़कर [नेटलिफाई के लिए कस्टम हेडर][ntl-headers] सेट कर सकते हैं:

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## शक्ति की खपत

एक वेब पेज कैसे बनाया जाता है यह उपयोगकर्ता के उपकरण पर चलने में लगने वाली शक्ति को प्रभावित कर सकता है।
न्यूनतम Javascript का उपयोग करके, Starlight उपयोगकर्ता के फोन, टैबलेट या कंप्यूटर को पेजो को लोड करने और प्रस्तुत करने के लिए आवश्यक प्रसंस्करण शक्ति की मात्रा को कम कर देता है।

वैश्लेषिकी ट्रैकिंग स्क्रिप्ट या वीडियो एम्बेड जैसी JavaScript-भारी सामग्री सुविधाएं जोड़ते समय सावधान रहें क्योंकि ये पेज शक्ति उपयोग को बढ़ा सकते हैं।
यदि आपको विश्लेषण की आवश्यकता है, तो [Cabin][cabin], [Fathom][fathom] या [Plausible][plausible] जैसे हल्के विकल्प चुनने पर विचार करें।
[उपयोगकर्ता इंटरैक्शन पर वीडियो लोड][lazy-video] होने की प्रतीक्षा करके YouTube और Vimeo वीडियो जैसे एंबेड को बेहतर बनाया जा सकता है।
[`astro-embed`][embed] जैसे पैकेज सामान्य सेवाओं के लिए मदद कर सकते हैं।

:::tip[क्या आप जानते हैं?]
JavaScript को पार्स करना और संकलित करना ब्राउज़र द्वारा किए जाने वाले सबसे महंगे कार्यों में से एक है।
समान आकार की JPEG छवि प्रस्तुत करने की तुलना में, [JavaScript को संसाधित होने में 30 गुना से अधिक समय लग सकता है][cost-of-js]।
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## होस्टिंग

जहां एक वेब पेज होस्ट किया गया है, उसका इस बात पर बड़ा प्रभाव पड़ सकता है कि आपकी दस्तावेज़ीकरण साइट पर्यावरण के अनुकूल कितनी है।
डेटा सेंटर और सर्वर फ़ार्म का बड़ा पारिस्थितिक प्रभाव हो सकता है, जिसमें उच्च बिजली की खपत और पानी का गहन उपयोग शामिल है।

नवीकरणीय शक्ति का उपयोग करने वाले होस्ट को चुनने का मतलब आपकी साइट के लिए कम कार्बन उत्सर्जन होगा। [Green Web Directory][gwb] एक उपकरण है जो आपको होस्टिंग कंपनियों को ढूंढने में मदद कर सकता है।

[gwb]: https://www.thegreenwebfoundation.org/directory/

## तुलना

क्या आप जानना चाहते हैं कि अन्य दस्तावेज़ीकरण फ्रेमवर्क की तुलना कैसे की जाती है?
[Website Carbon Calculator][wcc] के साथ ये परीक्षण विभिन्न उपकरणों के साथ बनाए गए समान पृष्ठों की तुलना करते हैं।

| फ्रेमवर्क                   | प्रति पृष्ठ विज़िट CO₂ |
| --------------------------- | ---------------------- |
| [Starlight][sl-carbon]      | 0.01g                  |
| [VitePress][vp-carbon]      | 0.05g                  |
| [Docus][dc-carbon]          | 0.05g                  |
| [Sphinx][sx-carbon]         | 0.07g                  |
| [MkDocs][mk-carbon]         | 0.10g                  |
| [Nextra][nx-carbon]         | 0.11g                  |
| [docsify][dy-carbon]        | 0.11g                  |
| [Docusaurus][ds-carbon]     | 0.24g                  |
| [Read the Docs][rtd-carbon] | 0.24g                  |
| [GitBook][gb-carbon]        | 0.71g                  |

<small>डेटा 14 मई 2023 को एकत्र किया गया। नवीनतम आंकड़े देखने के लिए लिंक पर क्लिक करें।</small>

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

## और अधिक संसाधन

### उपकरण

- [Website Carbon Calculator][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [WebPageTest Carbon Control](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### लेख और वार्ता

- [“Building a greener web”](https://youtu.be/EfPoOt7T5lg), Michelle Barker द्वारा वार्ता
- [“Sustainable Web Development Strategies Within An Organization”](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/), Michelle Barker का लेख
- [“A sustainable web for everyone”](https://2021.stateofthebrowser.com/speakers/tom-greenwood/), Tom Greenwood द्वारा वार्ता
- [“How Web Content Can Affect Power Usage”](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/), Benjamin Poulain और Simon Fraser का लेख

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/hi/guides/assets/
[islands]: https://docs.astro.build/hi/concepts/islands/
[wcc]: https://www.websitecarbon.com/
