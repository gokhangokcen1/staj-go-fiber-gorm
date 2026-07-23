import axios from "axios"

const API_URL = "http://localhost:3001/api"

export default {
  paketGonder(packetData) {
    return axios.post(`${API_URL}/packet-sender`, packetData)
  }
}
