<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()
const tools = [
  ['Subnetting', 'IP bloklarını ve kullanılabilir adres aralığını hesaplayın.', '/subnetting', '01'],
  ['Port Checker', 'Tek bir hedef portunun erişilebilirliğini kontrol edin.', '/portchecker', '02'],
  ['IP Scanner', 'Yerel ağdaki aktif cihazları hızlıca keşfedin.', '/ipscanner', '03'],
  ['SSL Checker', 'Sertifika zinciri ve güvenlik ayrıntılarını inceleyin.', '/sslcheck', '04'],
  ['WireHamsi', 'Ağ trafiğini canlı takip edip filtreleyin.', '/wirehamsi', '05'],
  ['Packet Sender', 'Özel TCP veya UDP paketleri oluşturup gönderin.', '/packet-sender', '06'],
  ['DNS Checker', 'DNS kayıtlarını ve A kaydı konumlarını doğrulayın.', '/dnschecker', '07'],
  ['WHOIS Sorgusu', 'Alan adı kayıt, kayıt firması ve iletişim bilgilerini inceleyin.', '/whois', '08'],
]
</script>

<template>
  <div class="home-page">
    <section class="hero-block">
      <p class="eyebrow">NETWORK TOOLKIT <span>● ONLINE</span></p>
      <h1>Ağınızı daha net<br /><em>görün.</em></h1>
      <p class="hero-copy">Tanılama, analiz ve test için ihtiyaç duyduğunuz araçlar sade, hızlı ve tek bir çalışma alanında.</p>
      <button @click="router.push('/ipscanner')">Taramaya başla <span>→</span></button>
    </section>
    <section class="tool-grid" aria-label="Ağ araçları">
      <button v-for="([title, description, path, number]) in tools" :key="path" class="tool-card" @click="router.push(path)">
        <span class="tool-number">{{ number }}</span><span class="tool-arrow">↗</span>
        <strong>{{ title }}</strong><small>{{ description }}</small>
      </button>
    </section>
  </div>
</template>

<style scoped>
.home-page { display: grid; gap: 26px; }
.hero-block { position: relative; overflow: hidden; padding: clamp(32px, 6vw, 76px); border-radius: 26px; border: 1px solid rgba(151, 223, 236, .16); background: linear-gradient(125deg, rgba(19, 53, 72, .9), rgba(13, 28, 52, .84)); box-shadow: 0 28px 75px rgba(0,0,0,.25); }
.hero-block::after { content: ''; position: absolute; width: 31rem; height: 31rem; right: -14rem; top: -17rem; border: 1px solid rgba(99, 234, 204, .28); border-radius: 50%; box-shadow: 0 0 0 42px rgba(99,234,204,.035), 0 0 0 84px rgba(99,234,204,.025); }
.eyebrow { margin: 0 0 19px; color: #81aabd; font-size: .7rem; letter-spacing: .17em; font-weight: 800; }.eyebrow span { color: #61e9c8; margin-left: 12px; }
h1 { position: relative; z-index: 1; margin: 0; font-size: clamp(2.45rem, 6vw, 5.2rem); line-height: .98; letter-spacing: -.07em; } h1 em { color: #65e5cc; font-style: normal; }
.hero-copy { position: relative; z-index: 1; max-width: 530px; margin: 25px 0; color: #bdd0dd; line-height: 1.8; font-size: 1.02rem; }
.hero-block > button { position: relative; z-index: 1; background: #64e4c9; border: 0; color: #06201f; font-weight: 800; padding: 14px 18px; border-radius: 10px; }.hero-block > button span { margin-left: 18px; font-size: 1.2rem; }
.tool-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(210px, 1fr)); gap: 12px; }
.tool-card { position: relative; min-height: 176px; padding: 22px; text-align: left; border-radius: 17px; border: 1px solid rgba(172,214,239,.12); background: rgba(18, 38, 61, .63); color: #eaf3ff; transition: transform .2s, border-color .2s, background .2s; }.tool-card:hover { transform: translateY(-4px); border-color: rgba(100,228,201,.55); background: rgba(28, 54, 76, .9); }.tool-card strong, .tool-card small { display: block; }.tool-card strong { margin-top: 42px; font-size: 1rem; }.tool-card small { margin-top: 9px; color: #9eb5c6; line-height: 1.55; font-size: .76rem; font-weight: 500; }.tool-number { color: #6ce3c9; font: 500 .72rem 'DM Mono', monospace; }.tool-arrow { position: absolute; right: 20px; top: 18px; color: #8fadc0; font-size: 1.2rem; }
</style>
