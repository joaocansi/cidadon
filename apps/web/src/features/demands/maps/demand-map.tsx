"use client";

import "leaflet/dist/leaflet.css";

import { useEffect, useMemo } from "react";
import { CircleMarker, MapContainer, TileLayer, useMap, useMapEvents } from "react-leaflet";

import { DEMAND_STATUS_META } from "@/features/demands/components/demand-status-badge";
import type { Demand } from "@/lib/api";
import { cn } from "@/lib/utils";

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

function DemandViewport({ demands, enabled }: { demands: Demand[]; enabled?: boolean }) {
  const map = useMap();
  const positions = useMemo(
    () =>
      demands
        .filter((demand) => Number.isFinite(demand.latitude) && Number.isFinite(demand.longitude))
        .map((demand) => [demand.latitude, demand.longitude] as [number, number]),
    [demands],
  );

  useEffect(() => {
    if (!enabled || positions.length === 0) return;
    if (positions.length === 1) {
      map.setView(positions[0], 15);
      return;
    }
    map.fitBounds(positions, { padding: [44, 44], maxZoom: 14 });
  }, [enabled, map, positions]);

  return null;
}

function DemandMapLegend() {
  return (
    <div
      aria-label="Legenda dos status das demandas"
      className="border-line bg-card/95 absolute top-14 left-3 z-10 rounded-xl border px-3 py-2 shadow-sm backdrop-blur"
    >
      <p className="text-ink-faint mb-1.5 text-[10px] font-semibold tracking-wide uppercase">
        Status das demandas
      </p>
      <ul className="grid gap-1.5 sm:grid-cols-2">
        {Object.entries(DEMAND_STATUS_META).map(([status, meta]) => (
          <li key={status} className="text-ink-soft flex items-center gap-1.5 text-xs font-medium">
            <span
              className="size-2.5 shrink-0 rounded-full ring-1 ring-black/10"
              style={{ backgroundColor: meta.markerColor }}
              aria-hidden="true"
            />
            {meta.label}
          </li>
        ))}
      </ul>
    </div>
  );
}

export function DemandMap({
  demands = [],
  selectedId,
  onSelect,
  location,
  onLocationChange,
  className,
  fitToDemands = false,
  showLegend = false,
}: {
  demands?: Demand[];
  selectedId?: number;
  onSelect?: (demand: Demand) => void;
  location?: [number, number];
  onLocationChange?: (value: [number, number]) => void;
  className?: string;
  fitToDemands?: boolean;
  showLegend?: boolean;
}) {
  const initial =
    location ??
    (demands[0] ? ([demands[0].latitude, demands[0].longitude] as [number, number]) : center);
  const mapClassName = cn("demand-map h-full w-full", className ? undefined : "rounded-[8px]");
  return (
    <div
      className={cn(
        "relative isolate overflow-hidden",
        className ?? "h-[460px] w-full rounded-[8px]",
      )}
    >
      <MapContainer center={initial} zoom={13} scrollWheelZoom className={mapClassName}>
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <DemandViewport demands={demands} enabled={fitToDemands} />
        {demands
          .filter((demand) => Number.isFinite(demand.latitude) && Number.isFinite(demand.longitude))
          .map((demand) => {
            const selected = selectedId === demand.id;
            return (
              <CircleMarker
                key={demand.id}
                center={[demand.latitude, demand.longitude]}
                radius={selected ? 12 : 8}
                pathOptions={{
                  color: selected ? "#1f3d2a" : "#fffef8",
                  fillColor: DEMAND_STATUS_META[demand.status].markerColor,
                  fillOpacity: 1,
                  weight: selected ? 4 : 2,
                }}
                eventHandlers={{ click: () => onSelect?.(demand) }}
              />
            );
          })}
        {onLocationChange ? <LocationPicker value={location} onChange={onLocationChange} /> : null}
      </MapContainer>
      {showLegend ? <DemandMapLegend /> : null}
    </div>
  );
}
