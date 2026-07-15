import axios from "axios"

const API_URL = "http://localhost:3001/api/sslcheck"

export default {
  sslcheck(website, port) {
    return axios.post(API_URL, { website, port })
  }
}