<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

const props = defineProps({
  results: { type: Array, default: () => [] },
  recordType: { type: String, default: 'A' },
  domain: { type: String, default: '' },
})

const mapElement = ref(null)
let map
let markers

function statusText(result) {
  if (result.status === 'found') return 'Kayit bulundu'
  if (result.status === 'not_found') return 'Bulunamadi'
  return 'Sorgulanamadi'
}

function updateMap() {
  if (!map) return
  markers.clearLayers()
  const points = []
  props.results.forEach((result) => {
    const point = [result.latitude, result.longitude]
    points.push(point)
    const color = result.status === 'found' ? '#22c55e' : result.status === 'not_found' ? '#f59e0b' : '#ef4444'
    const marker = L.circleMarker(point, { radius: 9, color, weight: 2, fillColor: color, fillOpacity: .85 })
    marker.bindPopup(`<strong>${result.server}</strong><br>${result.location}<br>${statusText(result)}`)
    markers.addLayer(marker)
  })
  if (points.length) map.fitBounds(points, { padding: [38, 38], maxZoom: 4 })
  else map.setView([25, 8], 2)
}

onMounted(async () => {
  await nextTick()
  map = L.map(mapElement.value, { scrollWheelZoom: false }).setView([25, 8], 2)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19, attribution: '&copy; OpenStreetMap contributors',
  }).addTo(map)
  markers = L.layerGroup().addTo(map)
  updateMap()
})

watch(() => props.results, updateMap, { deep: true })
onBeforeUnmount(() => map?.remove())
</script>

<template>
  <section class="dns-map-card">
    <h3>{{ recordType }} DNS Yayilim Haritasi</h3>
    <div ref="mapElement" class="map" aria-label="DNS sunuculari haritasi"></div>
    <p v-if="!results.length" class="map-empty">Sonuclari haritada gormek icin bir alan adi ve kayit turu ile kontrol edin.</p>

    <div v-if="results.length" class="resolver-results">
      <h4>{{ domain }} — {{ recordType }} sunucu cevaplari</h4>
      <article v-for="result in results" :key="result.server" class="resolver-result" :class="result.status">
        <div class="resolver-heading">
          <span><strong>{{ result.server }}</strong><small>{{ result.location }}</small></span>
          <b>{{ statusText(result) }}</b>
        </div>
        <ul v-if="result.answers?.length">
          <li v-for="answer in result.answers" :key="answer">{{ answer }}</li>
        </ul>
        <p v-else>{{ result.error || 'Bu sunucuda secilen kayit turu bulunamadi.' }}</p>
      </article>
    </div>
  </section>
</template>

<style scoped>
.dns-map-card { margin-top: 24px; padding: 22px; border: 1px solid rgba(255,255,255,.12); border-radius: 18px; background: rgba(255,255,255,.04); }
.map { height: 420px; overflow: hidden; border-radius: 14px; background: #102235; }
.map-empty { margin: 14px 0 0; color: #a9bed2; }
.resolver-results { margin-top: 18px; }
.resolver-results h4 { margin: 0 0 10px; }
.resolver-result { margin: 10px 0; padding: 13px 15px; border: 1px solid rgba(255,255,255,.12); border-left-width: 4px; border-radius: 10px; background: rgba(7,22,38,.55); }
.resolver-result.found { border-left-color: #22c55e; }.resolver-result.not_found { border-left-color: #f59e0b; }.resolver-result.error { border-left-color: #ef4444; }
.resolver-heading { display: flex; justify-content: space-between; gap: 12px; align-items: center; }.resolver-heading small { display: block; color: #a9bed2; margin-top: 2px; }.resolver-heading b { font-size: .85rem; white-space: nowrap; }
ul { margin: 10px 0 0; padding-left: 20px; } li { overflow-wrap: anywhere; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: .86rem; } p { margin: 10px 0 0; color: #c2cfdb; }
</style>
