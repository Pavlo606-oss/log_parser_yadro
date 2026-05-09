CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    file_type TEXT NOT NULL,
    imported_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE nodes (
    id BIGSERIAL PRIMARY KEY,
    log_id BIGINT NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    node_guid TEXT NOT NULL,
    node_desc TEXT NOT NULL,
    num_ports INTEGER NOT NULL,
    node_type INTEGER NOT NULL,
    class_version INTEGER NOT NULL,
    base_version INTEGER NOT NULL,
    system_image_guid TEXT NOT NULL,
    port_guid TEXT NOT NULL,
    UNIQUE (log_id, node_guid)
);

CREATE TABLE ports (
    id BIGSERIAL PRIMARY KEY,
    log_id BIGINT NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    node_id BIGINT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    node_guid TEXT NOT NULL,
    port_guid TEXT NOT NULL,
    port_num INTEGER NOT NULL,
    lid INTEGER NOT NULL,
    port_phy_state INTEGER NOT NULL,
    port_state INTEGER NOT NULL,
    link_width_actv INTEGER NOT NULL,
    link_speed_actv INTEGER NOT NULL,
    UNIQUE (log_id, node_guid, port_num)
);

CREATE TABLE switches (
    id BIGSERIAL PRIMARY KEY,
    log_id BIGINT NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    node_id BIGINT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    node_guid TEXT NOT NULL,
    linear_fdb_cap INTEGER NOT NULL,
    random_fdb_cap INTEGER NOT NULL,
    mcast_fdb_cap INTEGER NOT NULL,
    linear_fdb_top INTEGER NOT NULL,
    def_port INTEGER NOT NULL,
    def_mcast_pri_port INTEGER NOT NULL,
    def_mcast_not_pri_port INTEGER NOT NULL,
    life_time_value INTEGER NOT NULL,
    port_state_change INTEGER NOT NULL,
    optimized_slvl_mapping INTEGER NOT NULL,
    lids_per_port INTEGER NOT NULL,
    part_enf_cap INTEGER NOT NULL,
    inb_enf_cap INTEGER NOT NULL,
    outb_enf_cap INTEGER NOT NULL,
    filter_raw_inb_cap INTEGER NOT NULL,
    filter_raw_outb_cap INTEGER NOT NULL,
    enp0 INTEGER NOT NULL,
    mcast_fdb_top INTEGER NOT NULL,
    UNIQUE (log_id, node_guid)
);

CREATE TABLE nodes_info (
    id BIGSERIAL PRIMARY KEY,
    log_id BIGINT NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    node_id BIGINT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    node_guid TEXT NOT NULL,
    serial_number TEXT NOT NULL,
    part_number TEXT NOT NULL,
    revision TEXT NOT NULL,
    product_name TEXT NOT NULL,
    UNIQUE (log_id, node_guid)
);

CREATE TABLE sharp_an_info (
    id BIGSERIAL PRIMARY KEY,
    log_id BIGINT NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    sw_guid TEXT NOT NULL,
    endianness INTEGER NOT NULL,
    enable_endianness_per_job INTEGER NOT NULL,
    reproducibility_disable INTEGER NOT NULL,
    UNIQUE (log_id, sw_guid)
);