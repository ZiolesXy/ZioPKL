export type UserRole = "Admin" | "User"

export interface Address {
  uid: string
  label: string
  recipient_name: string
  phone: string
  address_line: string
  city: string
  province: string
  postal_code: string
  created_at: string
  updated_at: string
}

export interface User {
  id: number
  name: string
  email: string
  telephone_number?: string
  role: UserRole
  phone?: string
  profile_image_url?: string
  address?: Address[]
  created_at?: string
  updated_at?: string
}

export interface LoginData {
  user: User
  access_token: string
  refresh_token: string
}

export interface LoginResponse {
  status: string
  message: string
  data: LoginData
}

export interface ChangePasswordResponse {
  status: string
  message: string
}

export interface ChangePasswordData {
  old_password: string
  new_password: string
  confirm_password: string
}

export interface AddressResponse {
  status: string
  message: string
  data: AddressData[]
}

export interface AddressData {
  uid: string
  label: string
  recipient_name: string
  phone: string
  address_line: string
  city: string
  province: string
  postal_code: string
  is_primary: boolean
}

export type CreateAddressData = Omit<AddressData, "uid">