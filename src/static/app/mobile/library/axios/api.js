import axios from "axios"

const api = axios.create({
  baseURL: '/v1/api',
  headers: {
    'Content-Type': 'application/json'
  },
  transformRequest: [(data) => {
    if (!data) {
      return ''
    }
    return JSON.stringify(data.data)
  }]
})

api.interceptors.request.use(function (config) {
  const { limit, offset, page, unit } = { limit: 10, offset: 0, page: 0, unit: PAGINATION_UNIT, ...config.pagination}
  const start = page ? (limit * (page - 1)) : offset
  const end = page ? (limit * page - 1) : (offset + limit - 1)
  config.headers['Range'] = `${unit}=${start}-${end}`
  return config
})

api.interceptors.response.use(function (response) {
  const { data } = response
  if (data instanceof Array) {
    const pagination = paginate(response.headers['content-range'])
    const result = {
      pagination: pagination
    }
    result[pagination.unit] = data
    return result
  }
  return data
}, function (err) {
  err.message = err.response.data
  throw err
})

export default api

