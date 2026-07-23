<script setup>
import { RouterLink, RouterView } from 'vue-router'

const navigation = [
  ['/', 'Ana Sayfa'], ['/subnetting', 'Subnet'], ['/portchecker', 'Port'],
  ['/ipscanner', 'IP Tarama'], ['/sslcheck', 'SSL'], ['/wirehamsi', 'WireHamsi'],
  ['/packet-sender', 'Paket'], ['/dnschecker', 'DNS'],
  ['/whois', 'WHOIS'],
]
</script>

<template>
  <div class="app-shell">
    <header class="app-header">
      <div class="app-brand"><h1>Gökhan Gökcen</h1><p>Ağ araçları çalışma alanı</p></div>
      <nav class="nav-actions" aria-label="Ana navigasyon">
        <RouterLink v-for="([to, label]) in navigation" :key="to" :to="to">{{ label }}</RouterLink>
      </nav>
    </header>
    <main class="content-wrapper"><RouterView /></main>
  </div>
</template>
