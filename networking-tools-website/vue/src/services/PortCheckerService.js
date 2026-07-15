import axios from "axios"

const API_URL = "http://localhost:3001/api"

export default {
  kontrolEt(ip, port) {
    return axios.post(`${API_URL}/portcheck`, { ip, port })
  },
  mevcutIPGetir() {
    return axios.get(`${API_URL}/myip`)
  }
}
