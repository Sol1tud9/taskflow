import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api'

export default function Dashboard() {
  const [stats, setStats] = useState({
    tasks: { total: 0, todo: 0, in_progress: 0, done: 0 },
    users: 0,
    teams: 0,
    activities: 0,
  })
  const [recentTasks, setRecentTasks] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [])

  async function loadData() {
    try {
      const [tasksRes, activitiesRes] = await Promise.all([
        api.tasks.list({ limit: 5 }),
        api.activities.list({ limit: 10 }),
      ])
      
      setRecentTasks(tasksRes.tasks || [])
      setStats(prev => ({
        ...prev,
        tasks: {
          total: tasksRes.total || 0,
          todo: (tasksRes.tasks || []).filter(t => t.status === 'todo').length,
          in_progress: (tasksRes.tasks || []).filter(t => t.status === 'in_progress').length,
          done: (tasksRes.tasks || []).filter(t => t.status === 'done').length,
        },
        activities: activitiesRes.total || 0,
      }))
    } catch (error) {
      console.error('Failed to load dashboard data:', error)
    } finally {
      setLoading(false)
    }
  }

  const statCards = [
    { label: 'Total Tasks', value: stats.tasks.total, color: 'bg-primary-500', icon: TaskIcon },
    { label: 'To Do', value: stats.tasks.todo, color: 'bg-amber-500', icon: TodoIcon },
    { label: 'In Progress', value: stats.tasks.in_progress, color: 'bg-blue-500', icon: ProgressIcon },
    { label: 'Completed', value: stats.tasks.done, color: 'bg-emerald-500', icon: DoneIcon },
  ]

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-white">Dashboard</h1>
        <p className="text-dark-400 mt-1">Welcome to TaskFlow</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {statCards.map(({ label, value, color, icon: Icon }) => (
          <div key={label} className="card">
            <div className="flex items-center gap-4">
              <div className={`${color} p-3 rounded-xl`}>
                <Icon className="w-6 h-6 text-white" />
              </div>
              <div>
                <p className="text-2xl font-bold text-white">{loading ? '-' : value}</p>
                <p className="text-sm text-dark-400">{label}</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-white">Recent Tasks</h2>
            <Link to="/tasks" className="text-sm text-primary-400 hover:text-primary-300">
              View all
            </Link>
          </div>
          
          {loading ? (
            <div className="space-y-3">
              {[...Array(3)].map((_, i) => (
                <div key={i} className="h-16 bg-dark-800 rounded-lg animate-pulse" />
              ))}
            </div>
          ) : recentTasks.length === 0 ? (
            <div className="text-center py-8 text-dark-500">
              <TaskIcon className="w-12 h-12 mx-auto mb-2 opacity-50" />
              <p>No tasks yet</p>
              <Link to="/tasks" className="text-primary-400 hover:text-primary-300 text-sm">
                Create your first task
              </Link>
            </div>
          ) : (
            <div className="space-y-2">
              {recentTasks.map(task => (
                <div key={task.id} className="flex items-center gap-3 p-3 bg-dark-800/50 rounded-lg">
                  <StatusBadge status={task.status} />
                  <div className="flex-1 min-w-0">
                    <p className="text-dark-200 font-medium truncate">{task.title}</p>
                    <p className="text-xs text-dark-500">Priority: {task.priority}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-white">Quick Actions</h2>
          </div>
          
          <div className="grid grid-cols-2 gap-3">
            <Link
              to="/tasks"
              className="flex items-center gap-3 p-4 bg-dark-800/50 rounded-lg hover:bg-dark-800 transition-colors group"
            >
              <div className="p-2 bg-primary-500/10 rounded-lg group-hover:bg-primary-500/20 transition-colors">
                <TaskIcon className="w-5 h-5 text-primary-400" />
              </div>
              <span className="text-dark-200 font-medium">New Task</span>
            </Link>
            
            <Link
              to="/users"
              className="flex items-center gap-3 p-4 bg-dark-800/50 rounded-lg hover:bg-dark-800 transition-colors group"
            >
              <div className="p-2 bg-emerald-500/10 rounded-lg group-hover:bg-emerald-500/20 transition-colors">
                <UserIcon className="w-5 h-5 text-emerald-400" />
              </div>
              <span className="text-dark-200 font-medium">Add User</span>
            </Link>
            
            <Link
              to="/teams"
              className="flex items-center gap-3 p-4 bg-dark-800/50 rounded-lg hover:bg-dark-800 transition-colors group"
            >
              <div className="p-2 bg-violet-500/10 rounded-lg group-hover:bg-violet-500/20 transition-colors">
                <TeamIcon className="w-5 h-5 text-violet-400" />
              </div>
              <span className="text-dark-200 font-medium">Create Team</span>
            </Link>
            
            <Link
              to="/activities"
              className="flex items-center gap-3 p-4 bg-dark-800/50 rounded-lg hover:bg-dark-800 transition-colors group"
            >
              <div className="p-2 bg-amber-500/10 rounded-lg group-hover:bg-amber-500/20 transition-colors">
                <ActivityIcon className="w-5 h-5 text-amber-400" />
              </div>
              <span className="text-dark-200 font-medium">View Activity</span>
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

function StatusBadge({ status }) {
  const styles = {
    todo: 'bg-amber-500/10 text-amber-400',
    in_progress: 'bg-blue-500/10 text-blue-400',
    done: 'bg-emerald-500/10 text-emerald-400',
    cancelled: 'bg-dark-500/10 text-dark-400',
  }
  
  return (
    <span className={`badge ${styles[status] || styles.todo}`}>
      {status?.replace('_', ' ')}
    </span>
  )
}

function TaskIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
    </svg>
  )
}

function TodoIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  )
}

function ProgressIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
    </svg>
  )
}

function DoneIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
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

function TeamIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
    </svg>
  )
}

function ActivityIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
    </svg>
  )
}

