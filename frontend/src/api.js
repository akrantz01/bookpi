import axios from 'axios'

// Set constants
axios.defaults.baseURL = `${window.location.protocol}//${window.location.host}/api`
axios.defaults.withCredentials = true

// Send a request and handle errors gracefully
async function request (options) {
  try {
    return await axios(options)
  } catch (e) {
    if (!e.response) throw Error(`failed to send request: ${e.message}`)
    return e.response
  }
}

// Download a file
async function downloadFile (path, response) {
  // Decode payload if json to ensure not error
  if (response.status !== 200) {
    const parsed = JSON.parse(await response.data.text())
    if (parsed.status === 'error') return { status: response.status, reason: parsed.reason }
  }

  // Split to get file name
  const fileParts = path.split('/')

  // Create blob
  const url = window.URL.createObjectURL(new Blob([response.data]))

  // Generate and click link automatically
  const link = document.createElement('a')
  link.href = url
  link.setAttribute('download', fileParts[fileParts.length - 1])
  document.body.appendChild(link)
  link.click()

  // Cleanup
  link.remove()
  window.URL.revokeObjectURL(url)

  return { status: response.status }
}

const capitalize = str => str.charAt(0).toUpperCase() + str.slice(1)

class Authentication {
  static async register (username, password, name) {
    const response = await request({
      url: '/auth/register',
      method: 'post',
      data: { username, password, name },
      headers: { 'Content-Type': 'application/json' }
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async login (username, password) {
    const response = await request({
      url: '/auth/login',
      method: 'post',
      data: { username, password },
      headers: { 'Content-Type': 'application/json' }
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async logout () {
    const response = await request({
      url: '/auth/logout',
      method: 'get'
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }
}

class Users {
  static async read (username) {
    const response = await request({
      url: `/user/${username}`,
      method: 'get'
    })

    if (response.data.status === 'success') return { status: response.status, data: response.data.data }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async readSelf () {
    const response = await request({
      url: '/user',
      method: 'get'
    })

    if (response.data.status === 'success') return { status: response.status, data: response.data.data }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async update (name, password) {
    const response = await request({
      url: '/user',
      method: 'put',
      data: { name, password },
      headers: { 'Content-Type': 'application/json' }
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async delete () {
    const response = await request({
      url: '/user',
      method: 'delete'
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }
}

class Chats {
  static async list () {
    const response = await request({
      url: '/chats',
      method: 'get'
    })

    if (response.data.status === 'success') return { status: response.status, data: response.data.data }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async create (to, message) {
    const response = await request({
      url: '/chats',
      method: 'post',
      data: { to, message },
      headers: { 'Content-Type': 'application/json' }
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async read (id) {
    const response = await request({
      url: `/chats/${id}`,
      method: 'get'
    })

    if (response.data.status === 'success') return { status: response.status, data: response.data.data }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async delete (id) {
    const response = await request({
      url: `/chats/${id}`,
      method: 'delete'
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }
}

class Messages {
  static async list (id) {
    const response = await request({
      url: `/chats/${id}/messages`,
      method: 'get'
    })

    if (response.data.status === 'success') return { status: response.status, data: response.data.data }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async create (id, message) {
    const response = await request({
      url: `/chats/${id}/messages`,
      method: 'post',
      data: { message },
      headers: { 'Content-Type': 'application/json' }
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }
}

class Files {
  static async read (path, download = false) {
    const response = await request({
      url: `/files/${path}`,
      method: 'get',
      params: { download: (download) ? 'yes' : 'no' },
      responseType: (download) ? 'blob' : 'json'
    })

    if (download) return await downloadFile(path, response)

    if (response.data.status === 'success') return { status: response.status, data: response.data.data }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async create (path, file, directory = false, progressFunc = null) {
    const form = new FormData()
    form.append('file', file)

    const response = await request({
      url: `/files/${path}`,
      method: 'post',
      data: form,
      params: { directory: (directory) ? 'yes' : 'no' },
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: progressFunc
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async update (path, newFilename, newPath) {
    const response = await request({
      url: `/files/${path}`,
      method: 'put',
      data: { filename: newFilename, path: newPath },
      headers: { 'Content-Type': 'application/json' }
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async delete (path) {
    const response = await request({
      url: `/files/${path}`,
      method: 'delete'
    })

    if (response.data.status === 'success') return { status: response.status, data: response.data.data }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }
}

class Shares {
  static async list () {
    const response = await request({
      url: '/shares',
      method: 'get'
    })

    if (response.data.status === 'success') return { status: response.status, data: response.data.data }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async create (file, to) {
    const response = await request({
      url: '/shares',
      method: 'post',
      data: { file, to },
      headers: { 'Content-Type': 'application/json' }
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }

  static async download (user, file, describe = false) {
    const response = await request({
      url: `/shares/${user}/${file}`,
      method: 'get',
      params: { describe: (describe) ? 'yes' : 'no' },
      responseType: (describe) ? 'json' : 'blob'
    })

    if (describe) {
      if (response.data.status === 'success') return { status: response.status, data: response.data.data }
      return { status: response.status, reason: capitalize(response.data.reason) }
    }

    return await downloadFile(file, response)
  }

  static async delete (user, file) {
    const response = await request({
      url: `/shares/${user}/${file}`,
      method: 'delete'
    })

    if (response.data.status === 'success') return { status: response.status }
    return { status: response.status, reason: capitalize(response.data.reason) }
  }
}

export {
  Authentication,
  Users,
  Chats,
  Messages,
  Files,
  Shares
}
