import axios from "axios"

const API_URL = "http://localhost:3001/api"

export default {
  cihazlariGetir() {
    return axios.get(`${API_URL}/devices`)
  },
  baslat(device) {
    return axios.post(`${API_URL}/capture/start`, { device })
  },
  durdur() {
    return axios.post(`${API_URL}/capture/stop`)
  },
  websocketUrl() {
    return "ws://localhost:3001/ws/capture"
  }
}
