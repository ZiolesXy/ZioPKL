// types/flight.types.ts

import { Airline } from "./airline"
import { Airport } from "./airport"
import { FlightClassType } from "./base"

export interface FlightClass {
  id: number
  class_type: FlightClassType
  price: number
  total_seats: number
  seats?: number | null
}

export interface FlightSeat {
  id: number
  seat_number: string
  is_available: boolean
  class_type: FlightClassType
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
  meta: {
    limit: number;
    page: number;
    total: number;
  };
}

export interface Flight {
  id: number
  airline_id?: number
  origin_id?: number
  destination_id?: number
  flight_number: string
  departure_time: string
  arrival_time: string
  total_seats: number
  available_seats: number
  total_rows: number
  total_columns: number
  airline: Airline
  origin: Airport
  destination: Airport
  classes: FlightClass[]
}
