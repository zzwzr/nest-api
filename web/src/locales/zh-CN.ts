export default {
  common: {
    copy: '复制',
    logout: '退出登录',
    admin: '管理员',
    user: '用户',
  },
  topbar: {
    themeLight: '切换深色模式',
    themeDark: '切换浅色模式',
    settings: '系统设置',
    language: '语言',
  },
  admin: {
    title: '系统设置',
    desc: '超级管理员可查看全部用户与工作空间，并将工作空间转移给其他用户。',
    users: '用户列表',
    workspaces: '工作空间',
    transferTo: '转移给',
    selectUser: '选择用户',
    transferSuccess: '工作空间已转移',
    loadUsersFailed: '加载用户失败',
    loadWorkspacesFailed: '加载工作空间失败',
    transferFailed: '转移失败',
    columns: {
      id: 'ID',
      name: '名称',
      account: '账号',
      email: '邮箱',
      role: '角色',
      owner: '当前所有者',
      createdAt: '创建时间',
    },
  },
  home: {
    welcome: '欢迎回来，{name}',
    desc: '这里将展示你创建和加入的工作空间。接口开发页面会在后续版本中接入。',
  },
  auth: {
    logoutSuccess: '已退出登录',
  },
}
