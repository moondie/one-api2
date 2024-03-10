import {
  Typography,
  Stack,
  OutlinedInput,
  InputAdornment,
  Button,
  InputLabel,
  FormControl,
  FormControlLabel,
  MenuItem,
  Select,
  Checkbox
} from '@mui/material';
import { IconWallet } from '@tabler/icons-react';
import { useTheme } from '@mui/material/styles';
import SubCard from 'ui-component/cards/SubCard';
import UserCard from 'ui-component/cards/UserCard';

import { API } from 'utils/api';
import React, { useEffect, useState } from 'react';
import { showError, showInfo, showSuccess, renderQuota } from 'utils/common';
import {useSearchParams} from "react-router-dom";

const TopupCard = () => {
  const theme = useTheme();
  const [redemptionCode, setRedemptionCode] = useState(5);
  const [topUpLink, setTopUpLink] = useState('');
  const [userQuota, setUserQuota] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [checked, setChecked] = useState(true);
  const [payType, setPayType] = React.useState('wxpay');
  const [searchParams, setSearchParams] = useSearchParams();

  useEffect(() => {
    if (searchParams.get('trade_status') === 'TRADE_SUCCESS') {
      showSuccess('充值成功！');
    }
  }, [searchParams]);

  const handleChangePayType = (event) => {
    setPayType(event.target.value);
  };

  const topUp = async () => {
    if (redemptionCode === 0) {
      showInfo('请输入充值金额！');
      return;
    }
    setIsSubmitting(true);
    try {
      const res = await API.post('/api/user/recharge', {
        amount: redemptionCode,
        type: payType
      });
      const { success, message, payurl } = res.data;
      if (success) {
        showSuccess('创建充值链接成功！');
        setRedemptionCode(5);
        window.location.href = payurl;
      } else {
        showError(message);
      }
    } catch (err) {
      showError('请求失败');
    } finally {
      setIsSubmitting(false);
    }
  };

  const getUserQuota = async () => {
    try {
      let res = await API.get(`/api/user/self`);
      const { success, message, data } = res.data;
      if (success) {
        setUserQuota(data.quota);
      } else {
        showError(message);
      }
    } catch (error) {
      return;
    }
  };
  useEffect(() => {
    getUserQuota();
  }, []);

  return (
    <UserCard>
      <Stack direction="row" alignItems="center" justifyContent="center" spacing={2} paddingTop={'20px'}>
        <IconWallet color={theme.palette.primary.main} />
        <Typography variant="h4">当前额度:</Typography>
        <Typography variant="h4">{renderQuota(userQuota)}</Typography>
      </Stack>
      <SubCard
        sx={{
          marginTop: '40px'
        }}
      >
        <FormControl fullWidth variant="outlined">
          <InputLabel htmlFor="key">充值金额($)</InputLabel>
          <OutlinedInput
            id="key"
            label="充值金额($)"
            type="number"
            value={redemptionCode}
            onChange={(e) => {
              if (e.target.value < 0) {
                setRedemptionCode(1);
              } else if (e.target.value > 50) {
                setRedemptionCode(50);
              } else {
                setRedemptionCode(parseInt(e.target.value));
              }
            }}
            name="key"
            placeholder="请输入充值金额"
            endAdornment={
              <InputAdornment position="end">
                <Select sx={{ height: 30, marginRight: 1 }} value={payType} onChange={handleChangePayType}>
                  <MenuItem value={'wxpay'}>微信</MenuItem>
                  <MenuItem value={'alipay'}>支付宝</MenuItem>
                </Select>

                <Button variant="contained" onClick={topUp} disabled={isSubmitting}>
                  {isSubmitting ? '充值中...' : '充值'}
                </Button>
              </InputAdornment>
            }
            aria-describedby="helper-text-channel-quota-label"
          />
        </FormControl>

        {/*<Stack justifyContent="center" alignItems={'center'} spacing={3} paddingTop={'20px'}>*/}
        {/*  <Typography variant={'h4'} color={theme.palette.grey[700]}>*/}
        {/*    还没有兑换码？ 点击获取兑换码：*/}
        {/*  </Typography>*/}
        {/*  <Button variant="contained" onClick={openTopUpLink}>*/}
        {/*    获取兑换码*/}
        {/*  </Button>*/}
        {/*</Stack>*/}
      </SubCard>
    </UserCard>
  );
};

export default TopupCard;
