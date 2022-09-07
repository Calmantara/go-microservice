package balance

func (b *BalanceRepoImpl) GenerateWalletBalanceView() string {
	return `create or replace view balance_threshold as
			with agg_table as(
				select 
					lb.id, 
					lb.wallet_id, 
					lb.amount,
					first_value(lb.created_at) 
						over(
							w
							range between interval '2 minutes' preceding and current row  
						) as first_timestamp,
					sum(lb.amount) over(w range between current row and interval '2 minutes' following ) as agg_amount,
					case when 
						sum(lb.amount) over(w range between current row and interval '2 minutes' following ) > 10000 
							then true 
							else false 
					end as over_threshold,
					sum(lb.amount) 
						over w as total_amount
				from h_balances lb
				where lb.record_flag = 'ACTIVE'
				window w as (partition by lb.wallet_id order by lb.created_at asc)
			)
			select 
				foo.wallet_id, 
				foo.first_timestamp, 
				foo.total_balance,
				foo.total_per_window
			from( 
					select
						ag.*,
						last_value(ag.total_amount) over w as total_balance,
						first_value(ag.agg_amount) over (w range between current row and interval '2 minutes' following) as total_per_window
					from agg_table ag
					window w as (partition by ag.wallet_id order by ag.first_timestamp)
				) foo
			group by foo.wallet_id, foo.first_timestamp, foo.total_balance, foo.total_per_window
			order by wallet_id, first_timestamp;`
}

func (b *BalanceRepoImpl) GetWalletBalanceView() string {
	return `select 
				bt.wallet_id,
				bt.total_balance as amount,
				bt.total_per_window,
				case 
					when bt.total_per_window > 10000  and (bt.first_timestamp + interval '2 minutes') > ? 
						then true
					else false
				end as above_threshold
			from balance_threshold bt
			where bt.wallet_id = ?
			order by bt.first_timestamp desc
			limit 1;`
}

func (b *BalanceRepoImpl) GetSumWalletBalance() string {
	return `select 
				bt.wallet_id as wallet_id,
				sum(bt.amount) as amount
			from h_balances bt
			where bt.wallet_id = ?
			group by bt.wallet_id;`
}
