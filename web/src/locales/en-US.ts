export default {
  common: {
    copy: 'Copy',
    logout: 'Sign out',
    admin: 'Admin',
    user: 'User',
  },
  topbar: {
    themeLight: 'Switch to dark mode',
    themeDark: 'Switch to light mode',
    settings: 'System settings',
    language: 'Language',
  },
  admin: {
    title: 'System settings',
    desc: 'Super admins can view all users and workspaces, and transfer workspaces to other users.',
    users: 'Users',
    workspaces: 'Workspaces',
    transferTo: 'Transfer to',
    selectUser: 'Select user',
    transferSuccess: 'Workspace transferred',
    loadUsersFailed: 'Failed to load users',
    loadWorkspacesFailed: 'Failed to load workspaces',
    transferFailed: 'Transfer failed',
    columns: {
      id: 'ID',
      name: 'Name',
      account: 'Account',
      email: 'Email',
      role: 'Role',
      owner: 'Owner',
      createdAt: 'Created at',
    },
  },
  home: {
    welcome: 'Welcome back, {name}',
    desc: 'Your created and joined workspaces will appear here. API development pages are coming soon.',
  },
  auth: {
    logoutSuccess: 'Signed out',
  },
}
