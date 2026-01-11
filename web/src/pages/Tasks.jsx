import { useState, useEffect } from 'react'
import { api } from '../api'
import Modal from '../components/Modal'

const STATUSES = ['todo', 'in_progress', 'done', 'cancelled']
const PRIORITIES = ['low', 'medium', 'high']

export default function Tasks() {
  const [tasks, setTasks] = useState([])
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [filter, setFilter] = useState({ status: '' })
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    priority: 'medium',
    creator_id: '',
    assignee_id: '',
    team_id: '',
  })

  useEffect(() => {
    loadTasks()
    loadUsers()
  }, [filter])

  async function loadUsers() {
    try {
      const res = await api.users.list()
      setUsers(Array.isArray(res) ? res : res.users || [])
    } catch (error) {
      console.error('Failed to load users:', error)
    }
  }

  async function loadTasks() {
    try {
      setLoading(true)
      const params = {}
      if (filter.status) params.status = filter.status
      const res = await api.tasks.list(params)
      setTasks(res.tasks || [])
    } catch (error) {
      console.error('Failed to load tasks:', error)
    } finally {
      setLoading(false)
    }
  }

  async function handleCreate(e) {
    e.preventDefault()
    try {
      await api.tasks.create(formData)
      setIsModalOpen(false)
      setFormData({ title: '', description: '', priority: 'medium', creator_id: '', assignee_id: '', team_id: '' })
      loadTasks()
    } catch (error) {
      console.error('Failed to create task:', error)
    }
  }

  async function handleStatusChange(id, status) {
    try {
      await api.tasks.update(id, { status })
      loadTasks()
    } catch (error) {
      console.error('Failed to update task:', error)
    }
  }

  async function handleDelete(id) {
    if (!confirm('Are you sure you want to delete this task?')) return
    try {
      await api.tasks.delete(id)
      loadTasks()
    } catch (error) {
      console.error('Failed to delete task:', error)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white">Tasks</h1>
          <p className="text-dark-400 mt-1">Manage your tasks</p>
        </div>
        <button onClick={() => setIsModalOpen(true)} className="btn btn-primary">
          <PlusIcon className="w-5 h-5" />
          New Task
        </button>
      </div>

      <div className="flex gap-2">
        <button
          onClick={() => setFilter({ status: '' })}
          className={`btn ${!filter.status ? 'btn-primary' : 'btn-secondary'}`}
        >
          All
        </button>
        {STATUSES.map(status => (
          <button
            key={status}
            onClick={() => setFilter({ status })}
            className={`btn ${filter.status === status ? 'btn-primary' : 'btn-secondary'}`}
          >
            {status.replace('_', ' ')}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="h-20 bg-dark-800 rounded-xl animate-pulse" />
          ))}
        </div>
      ) : tasks.length === 0 ? (
        <div className="card text-center py-12">
          <TaskIcon className="w-16 h-16 mx-auto text-dark-600 mb-4" />
          <h3 className="text-xl font-semibold text-dark-300 mb-2">No tasks found</h3>
          <p className="text-dark-500 mb-4">Create your first task to get started</p>
          <button onClick={() => setIsModalOpen(true)} className="btn btn-primary">
            <PlusIcon className="w-5 h-5" />
            Create Task
          </button>
        </div>
      ) : (
        <div className="space-y-3">
          {tasks.map(task => (
            <div key={task.id} className="card flex items-start gap-4">
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-1">
                  <h3 className="text-lg font-semibold text-white">{task.title}</h3>
                  <StatusBadge status={task.status} />
                  <PriorityBadge priority={task.priority} />
                </div>
                {task.description && (
                  <p className="text-dark-400 text-sm mb-2">{task.description}</p>
                )}
                <p className="text-xs text-dark-500">
                  ID: {task.id?.slice(0, 8)}...
                </p>
              </div>
              
              <div className="flex items-center gap-2">
                <select
                  value={task.status}
                  onChange={(e) => handleStatusChange(task.id, e.target.value)}
                  className="input py-1.5 text-sm w-32"
                >
                  {STATUSES.map(s => (
                    <option key={s} value={s}>{s.replace('_', ' ')}</option>
                  ))}
                </select>
                <button
                  onClick={() => handleDelete(task.id)}
                  className="btn btn-danger p-2"
                >
                  <TrashIcon className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title="Create Task">
        <form onSubmit={handleCreate} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-dark-300 mb-1">Title</label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              className="input"
              placeholder="Task title"
              required
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-dark-300 mb-1">Description</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="input min-h-[80px] resize-none"
              placeholder="Task description"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-dark-300 mb-1">Priority</label>
            <select
              value={formData.priority}
              onChange={(e) => setFormData({ ...formData, priority: e.target.value })}
              className="input"
            >
              {PRIORITIES.map(p => (
                <option key={p} value={p}>{p}</option>
              ))}
            </select>
          </div>
          
          <div>
            <label className="block text-sm font-medium text-dark-300 mb-1">Creator</label>
            {users.length > 0 ? (
              <select
                value={formData.creator_id}
                onChange={(e) => setFormData({ ...formData, creator_id: e.target.value })}
                className="input"
                required
              >
                <option value="">Select creator...</option>
                {users.map(user => (
                  <option key={user.id} value={user.id}>
                    {user.name} ({user.email})
                  </option>
                ))}
              </select>
            ) : (
              <p className="text-dark-400 text-sm">No users available. Create a user first.</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-dark-300 mb-1">Assignee (optional)</label>
            <select
              value={formData.assignee_id}
              onChange={(e) => setFormData({ ...formData, assignee_id: e.target.value })}
              className="input"
            >
              <option value="">No assignee</option>
              {users.map(user => (
                <option key={user.id} value={user.id}>
                  {user.name} ({user.email})
                </option>
              ))}
            </select>
          </div>
          
          <div className="flex gap-3 pt-2">
            <button type="button" onClick={() => setIsModalOpen(false)} className="btn btn-secondary flex-1">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary flex-1">
              Create
            </button>
          </div>
        </form>
      </Modal>
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
  return <span className={`badge ${styles[status] || styles.todo}`}>{status?.replace('_', ' ')}</span>
}

function PriorityBadge({ priority }) {
  const styles = {
    low: 'bg-dark-500/10 text-dark-400',
    medium: 'bg-amber-500/10 text-amber-400',
    high: 'bg-red-500/10 text-red-400',
  }
  return <span className={`badge ${styles[priority] || styles.medium}`}>{priority}</span>
}

function TaskIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
    </svg>
  )
}

function PlusIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
    </svg>
  )
}

function TrashIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
    </svg>
  )
}

