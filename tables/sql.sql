-- t_order_info
CREATE TABLE `t_order_info`
(
    `id`                 BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '',
    `order_id`           VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `alg_id`             SMALLINT            NOT NULL DEFAULT '0' COMMENT '',
    `sub_alg_id`         SMALLINT            NOT NULL DEFAULT '0' COMMENT '',
    `payload`            VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `token_id`           VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '支付币种',
    `amount_token`       DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT '支付币种金额',
    `amount_dp`          DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT '充值dp数量',
    `payment_hash`       VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `payment_status`     SMALLINT            NOT NULL DEFAULT '0' COMMENT '0-unpaid 1-paid',
    `order_status`       SMALLINT            NOT NULL DEFAULT '0' COMMENT '0-default 1-success 2-fail',
    `dp_status`          SMALLINT            NOT NULL DEFAULT '0' COMMENT '0-default 1-wait 2-ing 3-ok',
    `hedging_status`     SMALLINT            NOT NULL DEFAULT '0' COMMENT '0-default 1-wait 2-ing 3-ok',
    `timestamp`          BIGINT              NOT NULL DEFAULT '0' COMMENT '',
    `premium_percentage` DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT 'stripe溢价',
    `premium_base`       DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT 'stripe溢价',
    `premium_amount`     DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT 'stripe溢价',
    `created_at`         TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
    `updated_at`         TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_id` (`order_id`),
    KEY `k_timestamp` (`timestamp`),
    KEY `k_payload` (`payload`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='dp order info';


-- t_withdraw_info
CREATE TABLE `t_withdraw_info`
(
    `id`              BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '',
    `withdraw_id`     VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `order_id`        VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `alg_id`          SMALLINT            NOT NULL DEFAULT '0' COMMENT '',
    `sub_alg_id`      SMALLINT            NOT NULL DEFAULT '0' COMMENT '',
    `payload`         VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `token_id`        VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `amount_dp`       DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT '',
    `amount_token`    DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT '',
    `dp_hash`         VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `dp_status`       SMALLINT            NOT NULL DEFAULT '0' COMMENT '',
    `withdraw_hash`   VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `withdraw_status` SMALLINT            NOT NULL DEFAULT '0' COMMENT '',
    `timestamp`       BIGINT              NOT NULL DEFAULT '0' COMMENT '',
    `created_at`      TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
    `updated_at`      TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_withdraw_id` (`withdraw_id`),
    KEY `k_timestamp` (`timestamp`),
    KEY `k_payload` (`payload`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='dp withdraw info';

-- t_record_info
CREATE TABLE `t_record_info`
(
    `id`           BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '',
    `record_id`    VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `action`       VARCHAR(255)        NOT NULL DEFAULT '' COMMENT 'deposit, withdraw',
    `order_id`     VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `withdraw_id`  VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `alg_id`       SMALLINT            NOT NULL DEFAULT '0' COMMENT '',
    `sub_alg_id`   SMALLINT            NOT NULL DEFAULT '0' COMMENT '',
    `payload`      VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `token_id`     VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `amount_dp`    DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT '',
    `amount_token` DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT '',
    `timestamp`    BIGINT              NOT NULL DEFAULT '0' COMMENT '',
    `created_at`   TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
    `updated_at`   TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_wid` (`record_id`),
    KEY `k_timestamp` (`timestamp`),
    KEY `k_payload` (`payload`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='dp record info';

-- t_hedging_info
CREATE TABLE `t_hedging_info`
(
    `id`             BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '',
    `hedging_id`     VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `action`         VARCHAR(255)        NOT NULL DEFAULT '' COMMENT 'deposit, withdraw',
    `order_id`       VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `withdraw_id`    VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `token_id_from`  VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `token_id_to`    VARCHAR(255)        NOT NULL DEFAULT '' COMMENT '',
    `amount_from`    DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT '',
    `amount_to`      DECIMAL(50, 0)      NOT NULL DEFAULT '0' COMMENT '',
    `hedging_status` SMALLINT            NOT NULL DEFAULT '0' COMMENT '',
    `timestamp`      BIGINT              NOT NULL DEFAULT '0' COMMENT '',
    `created_at`     TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
    `updated_at`     TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_wid` (`hedging_id`),
    KEY `k_timestamp` (`timestamp`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='dp hedging info';