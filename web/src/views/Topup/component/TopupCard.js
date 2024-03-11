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
  Checkbox,
  Dialog,
  DialogTitle,
  Divider,
  DialogContent,
  DialogActions
} from '@mui/material';
import { IconWallet } from '@tabler/icons-react';
import { useTheme } from '@mui/material/styles';
import SubCard from 'ui-component/cards/SubCard';
import UserCard from 'ui-component/cards/UserCard';
import QRCode from 'qrcode.react';

import { API } from 'utils/api';
import React, { useEffect, useState } from 'react';
import { showError, showInfo, showSuccess, renderQuota } from 'utils/common';
import { useSearchParams } from 'react-router-dom';

const QRModal = ({ open, QRString, onCancel, onOk, quota }) => {
  return (
    <Dialog open={open} onClose={onCancel} sx={{ textAlign: 'center' }}>
      <DialogTitle sx={{ fontWeight: 700, lineHeight: '1.55556', padding: '24px', fontSize: '1.125rem' }}>
        微信支付二维码: 请支付:{quota}美元
      </DialogTitle>
      <Divider />
      <DialogContent sx={{width:'400px'}}>
        <QRCode value={QRString} />
        <DialogActions>
          <Button onClick={onCancel}>取消</Button>
          <Button onClick={onOk} type="submit" variant="contained" color="primary">
            已支付
          </Button>
        </DialogActions>
      </DialogContent>
    </Dialog>
  );
};

const TopupCard = () => {
  const theme = useTheme();
  const [redemptionCode, setRedemptionCode] = useState(5);
  const [userQuota, setUserQuota] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [payType, setPayType] = React.useState('wxpay');
  const [searchParams, setSearchParams] = useSearchParams();
  const [showModal, setShowModal] = useState(false);
  const [qrcode, setQrcode] = useState('');

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
        if (payurl.startsWith('weixin://')) {
          setQrcode(payurl);
          setShowModal(true);
          window.location.href = payurl;
        } else {
          window.location.href = payurl;
        }
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
    <>
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
      <QRModal
        open={showModal}
        onCancel={() => {
          setShowModal(false);
        }}
        onOk={() => {
          window.location.reload();
        }}
        QRString={qrcode}
        quota={redemptionCode}
      />
    </>
  );
};

export default TopupCard;
