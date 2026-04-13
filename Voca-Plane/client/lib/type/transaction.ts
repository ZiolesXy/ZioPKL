// types/transaction.ts

import { PaymentStatus, FlightClassType } from "./base"
import { Flight } from "./flight"
import { User } from "./user"

export interface TransactionPassenger {
  passenger_name: string
  nationality: string
  passport_no: string
  seat_number: string
  class_name: string
  price: number
}

export interface Transaction {
  id: number
  code: string
  user?: Pick<User, "id" | "name" | "email">
  flight: Flight
  transactions_passangers: TransactionPassenger[]
  total_price: number
  discount: number
  payment_url: string
  promo_code?: string
  payment_status: PaymentStatus
  created_at: string
  expires_at: string
}