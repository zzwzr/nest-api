export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface InstallStatus {
  installed: boolean
}

export interface DatabaseConfig {
  driver: 'postgres' | 'mysql'
  host: string
  port: number
  name: string
  user: string
  password: string
  ssl_mode: string
}

export interface AdminConfig {
  username: string
  password: string
  confirm_password: string
}

export interface InstallPayload {
  database: DatabaseConfig
  admin: AdminConfig
}

export interface TestDatabaseResult {
  ok: boolean
  message: string
}

export interface InstallResult {
  message: string
  database_user?: string
  database_password?: string
}
