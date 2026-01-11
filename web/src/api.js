const API_BASE = '/api/v1'

async function request(endpoint, options = {}) {
  const url = `${API_BASE}${endpoint}`
  const config = {
    headers: {
      'Content-Type': 'application/json',
    },
    ...options,
  }

  try {
    const response = await fetch(url, config)
    const data = await response.json()
    
    if (!response.ok) {
      throw new Error(data.error || 'Request failed')
    }
    
    return data
  } catch (error) {
    console.error('API Error:', error)
    throw error
  }
}

export const api = {
  users: {
    list: () => request('/users'),
    get: (id) => request(`/users/${id}`),
    create: (data) => request('/users', { method: 'POST', body: JSON.stringify(data) }),
    update: (id, data) => request(`/users/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
    activities: (id) => request(`/users/${id}/activities`),
  },
  
  teams: {
    list: () => request('/teams'),
    get: (id) => request(`/teams/${id}`),
    create: (data) => request('/teams', { method: 'POST', body: JSON.stringify(data) }),
    getMembers: (id) => request(`/teams/${id}/members`),
    addMember: (id, data) => request(`/teams/${id}/members`, { method: 'POST', body: JSON.stringify(data) }),
  },
  
  tasks: {
    list: (params = {}) => {
      const query = new URLSearchParams(params).toString()
      return request(`/tasks${query ? `?${query}` : ''}`)
    },
    get: (id) => request(`/tasks/${id}`),
    create: (data) => request('/tasks', { method: 'POST', body: JSON.stringify(data) }),
    update: (id, data) => request(`/tasks/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
    delete: (id) => request(`/tasks/${id}`, { method: 'DELETE' }),
    history: (id) => request(`/tasks/${id}/history`),
  },
  
  activities: {
    list: (params = {}) => {
      const query = new URLSearchParams(params).toString()
      return request(`/activities${query ? `?${query}` : ''}`)
    },
  },
}

