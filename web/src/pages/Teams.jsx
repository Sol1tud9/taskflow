import { useState, useEffect } from 'react'
import { api } from '../api'
import Modal from '../components/Modal'

export default function Teams() {
  const [teams, setTeams] = useState([])
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [isMembersModalOpen, setIsMembersModalOpen] = useState(false)
  const [isAddMemberModalOpen, setIsAddMemberModalOpen] = useState(false)
  const [selectedTeam, setSelectedTeam] = useState(null)
  const [members, setMembers] = useState([])
  const [formData, setFormData] = useState({ name: '', owner_id: '' })
  const [memberFormData, setMemberFormData] = useState({ user_id: '', role: 'member' })

  useEffect(() => {
    loadTeams()
    loadUsers()
  }, [])

  async function loadTeams() {
    try {
      setLoading(true)
      const res = await api.teams.list()
      setTeams(Array.isArray(res) ? res : res.teams || [])
    } catch (error) {
      console.error('Failed to load teams:', error)
      setTeams([])
    } finally {
      setLoading(false)
    }
  }

  async function loadUsers() {
    try {
      const res = await api.users.list()
      setUsers(Array.isArray(res) ? res : res.users || [])
    } catch (error) {
      console.error('Failed to load users:', error)
    }
  }

  async function handleCreate(e) {
    e.preventDefault()
    try {
      await api.teams.create(formData)
      setIsModalOpen(false)
      setFormData({ name: '', owner_id: '' })
      loadTeams()
    } catch (error) {
      console.error('Failed to create team:', error)
    }
  }

  async function viewMembers(team) {
    setSelectedTeam(team)
    try {
      const res = await api.teams.getMembers(team.id)
      setMembers(Array.isArray(res) ? res : res.members || [])
    } catch (error) {
      console.error('Failed to load members:', error)
      setMembers([])
    }
    setIsMembersModalOpen(true)
  }

  async function handleAddMember(e) {
    e.preventDefault()
    try {
      await api.teams.addMember(selectedTeam.id, memberFormData)
      setIsAddMemberModalOpen(false)
      setMemberFormData({ user_id: '', role: 'member' })
      const res = await api.teams.getMembers(selectedTeam.id)
      setMembers(Array.isArray(res) ? res : res.members || [])
    } catch (error) {
      console.error('Failed to add member:', error)
    }
  }

  function getUserName(userId) {
    const user = users.find(u => u.id === userId)
    return user ? user.name : userId
  }

  function getAvailableUsers() {
    const memberIds = members.map(m => m.user_id)
    return users.filter(u => !memberIds.includes(u.id))
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white">Teams</h1>
          <p className="text-dark-400 mt-1">Organize your users into teams</p>
        </div>
        <button onClick={() => setIsModalOpen(true)} className="btn btn-primary">
          <PlusIcon className="w-5 h-5" />
          Create Team
        </button>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="h-40 bg-dark-800 rounded-xl animate-pulse" />
          ))}
        </div>
      ) : teams.length === 0 ? (
        <div className="card text-center py-12">
          <TeamIcon className="w-16 h-16 mx-auto text-dark-600 mb-4" />
          <h3 className="text-xl font-semibold text-dark-300 mb-2">No teams yet</h3>
          <p className="text-dark-500 mb-4">Create your first team to collaborate</p>
          <button onClick={() => setIsModalOpen(true)} className="btn btn-primary">
            <PlusIcon className="w-5 h-5" />
            Create Team
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {teams.map(team => (
            <div key={team.id} className="card">
              <div className="flex items-start gap-4">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-violet-400 to-violet-600 flex items-center justify-center">
                  <TeamIcon className="w-6 h-6 text-white" />
                </div>
                <div className="flex-1">
                  <h3 className="text-lg font-semibold text-white">{team.name}</h3>
                  <p className="text-sm text-dark-400">Owner: {getUserName(team.owner_id)}</p>
                </div>
              </div>
              <div className="mt-4 pt-4 border-t border-dark-800 flex items-center justify-between">
                <p className="text-xs text-dark-500 font-mono break-all">
                  ID: {team.id}
                </p>
                <button 
                  onClick={() => viewMembers(team)}
                  className="text-sm text-primary-400 hover:text-primary-300"
                >
                  View Members
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title="Create Team">
        <form onSubmit={handleCreate} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-dark-300 mb-1">Team Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="input"
              placeholder="Engineering"
              required
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-dark-300 mb-1">Owner</label>
            {users.length > 0 ? (
              <select
                value={formData.owner_id}
                onChange={(e) => setFormData({ ...formData, owner_id: e.target.value })}
                className="input"
                required
              >
                <option value="">Select owner...</option>
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
          
          <div className="flex gap-3 pt-2">
            <button type="button" onClick={() => setIsModalOpen(false)} className="btn btn-secondary flex-1">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary flex-1" disabled={users.length === 0}>
              Create
            </button>
          </div>
        </form>
      </Modal>

      <Modal 
        isOpen={isMembersModalOpen} 
        onClose={() => setIsMembersModalOpen(false)} 
        title={`Members of ${selectedTeam?.name || 'Team'}`}
      >
        <div className="space-y-3">
          {members.length === 0 ? (
            <p className="text-dark-400 text-center py-4">No members in this team yet</p>
          ) : (
            members.map(member => (
              <div key={member.id} className="flex items-center gap-3 p-3 bg-dark-800 rounded-lg">
                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white font-bold">
                  {getUserName(member.user_id)?.charAt(0)?.toUpperCase() || 'U'}
                </div>
                <div className="flex-1">
                  <p className="text-white font-medium">{getUserName(member.user_id)}</p>
                  <p className="text-sm text-dark-400 capitalize">{member.role}</p>
                </div>
              </div>
            ))
          )}
        </div>
        <div className="mt-4 pt-4 border-t border-dark-800 flex gap-3">
          <button 
            onClick={() => {
              setIsAddMemberModalOpen(true)
            }} 
            className="btn btn-primary flex-1"
          >
            <PlusIcon className="w-4 h-4" />
            Add Member
          </button>
          <button 
            onClick={() => setIsMembersModalOpen(false)} 
            className="btn btn-secondary flex-1"
          >
            Close
          </button>
        </div>
      </Modal>

      <Modal
        isOpen={isAddMemberModalOpen}
        onClose={() => setIsAddMemberModalOpen(false)}
        title={`Add Member to ${selectedTeam?.name || 'Team'}`}
      >
        <form onSubmit={handleAddMember} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-dark-300 mb-1">User</label>
            {getAvailableUsers().length > 0 ? (
              <select
                value={memberFormData.user_id}
                onChange={(e) => setMemberFormData({ ...memberFormData, user_id: e.target.value })}
                className="input"
                required
              >
                <option value="">Select user...</option>
                {getAvailableUsers().map(user => (
                  <option key={user.id} value={user.id}>
                    {user.name} ({user.email})
                  </option>
                ))}
              </select>
            ) : (
              <p className="text-dark-400 text-sm">No available users to add.</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-dark-300 mb-1">Role</label>
            <select
              value={memberFormData.role}
              onChange={(e) => setMemberFormData({ ...memberFormData, role: e.target.value })}
              className="input"
            >
              <option value="member">Member</option>
              <option value="admin">Admin</option>
              <option value="viewer">Viewer</option>
            </select>
          </div>

          <div className="flex gap-3 pt-2">
            <button type="button" onClick={() => setIsAddMemberModalOpen(false)} className="btn btn-secondary flex-1">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary flex-1" disabled={getAvailableUsers().length === 0}>
              Add
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}

function TeamIcon({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
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
