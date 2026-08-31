"use client";

import "leaflet/dist/leaflet.css";

import { CircleMarker, MapContainer, TileLayer, useMapEvents } from "react-leaflet";

import type { Demand } from "@/lib/api";

const center: [number, number] = [-23.5505, -46.6333];

function LocationPicker({
  value,
  onChange,
}: {
  value?: [number, number];
  onChange?: (value: [number, number]) => void;
}) {
  useMapEvents({
    click(event) {
      onChange?.([event.latlng.lat, event.latlng.lng]);
    },
  });
  return value ? (
    <CircleMarker
      center={value}
      radius={10}
      pathOptions={{ color: "#355b3e", fillColor: "#b7d84b", fillOpacity: 1, weight: 3 }}
    />
  ) : null;
}

export function DemandMap({
  demands = [],
  selectedId,
  onSelect,
  location,
  onLocationChange,
  className,
}: {
  demands?: Demand[];
  selectedId?: number;
  onSelect?: (demand: Demand) => void;
  location?: [number, number];
  onLocationChange?: (value: [number, number]) => void;
  className?: string;
}) {
  const initial =
    location ??
    (demands[0] ? ([demands[0].latitude, demands[0].longitude] as [number, number]) : center);
  const mapClassName = `${className ?? "h-[460px] w-full rounded-[8px]"} demand-map`;
  return (
    <MapContainer center={initial} zoom={13} scrollWheelZoom className={mapClassName}>
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      {demands
        .filter((item) => item.latitude && item.longitude)
        .map((demand) => (
          <CircleMarker
            key={demand.id}
            center={[demand.latitude, demand.longitude]}
            radius={selectedId === demand.id ? 11 : 8}
            pathOptions={{
              color: "#1f3d2a",
              fillColor: selectedId === demand.id ? "#b7d84b" : "#e9a23b",
              fillOpacity: 1,
              weight: 2,
            }}
            eventHandlers={{ click: () => onSelect?.(demand) }}
          />
        ))}
      {onLocationChange ? <LocationPicker value={location} onChange={onLocationChange} /> : null}
    </MapContainer>
  );
}
