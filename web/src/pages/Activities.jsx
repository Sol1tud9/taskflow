import { useState, useEffect } from 'react'
import { api } from '../api'

export default function Activities() {
  const [activities, setActivities] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadActivities()
  }, [])

  async function loadActivities() {
    try {
      setLoading(true)
      const res = await api.activities.list({ limit: 50 })
      setActivities(res.activities || [])
    } catch (error) {
      console.error('Failed to load activities:', error)
      setActivities([])
    } finally {
      setLoading(false)
    }
  }

  function formatDate(dateStr) {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return '-'
    return date.toLocaleString()
  }

  function getActionColor(action) {
    switch (action) {
      case 'created': return 'text-emerald-400 bg-emerald-500/10'
      case 'updated': return 'text-blue-400 bg-blue-500/10'
      case 'deleted': return 'text-red-400 bg-red-500/10'
      default: return 'text-dark-400 bg-dark-500/10'
    }
  }

  function getEntityIcon(entityType) {
    switch (entityType) {
      case 'user': return <UserIcon className="w-4 h-4" />
      case 'task': return <TaskIcon className="w-4 h-4" />
      case 'team': return <TeamIcon className="w-4 h-4" />
      default: return <ActivityIcon className="w-4 h-4" />
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-white">Activities</h1>
        <p className="text-dark-400 mt-1">Track all activity in your organization</p>
      </div>

      {loading ? (
        <div className="space-y-3">
          {[...Array(10)].map((_, i) => (
            <div key={i} className="h-16 bg-dark-800 rounded-xl animate-pulse" />
          ))}
        </div>
      ) : activities.length === 0 ? (
        <div className="card text-center py-12">
          <ActivityIcon className="w-16 h-16 mx-auto text-dark-600 mb-4" />
          <h3 className="text-xl font-semibold text-dark-300 mb-2">No activities yet</h3>
          <p className="text-dark-500">Activities will appear here when users perform actions</p>
        </div>
      ) : (
        <div className="space-y-2">
          {activities.map((activity, index) => (
            <div key={activity.id || index} className="card py-4 flex items-center gap-4">
              <div className={`p-2 rounded-lg ${getActionColor(activity.action)}`}>
                {getEntityIcon(activity.entity_type)}
              </div>
              
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span className="text-dark-200 font-medium">
                    {activity.entity_type}
                  </span>
                  <span className={`badge ${getActionColor(activity.action)}`}>
                    {activity.action}
                  </span>
                </div>
                <p className="text-sm text-dark-500 break-all">
                  Entity: <span className="font-mono">{activity.entity_id}</span> | User: <span className="font-mono">{activity.user_id}</span>
                </p>
              </div>
              
              <div className="text-right">
                <p className="text-sm text-dark-400">{formatDate(activity.created_at || activity.createdAt)}</p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

function ActivityIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
    </svg>
  )
}

function UserIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
    </svg>
  )
}

function TaskIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
    </svg>
  )
}

function TeamIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
    </svg>
  )
}

