#include<bits/stdc++.h>
using namespace std;
#define maxq 300010
int n, k, N;
int a[maxq*2], sum[maxq*2];
struct Q
{
    int l, r, k;
    void input(){
        scanf("%d%d%d", &l, &r, &k);
        r++;
        k%=2;
    }
}q[maxq];
vector<int> v;
int seg[maxq<<3], num[maxq<<3];
bool tg[2][maxq<<3];
void push_up(int o){seg[o]=seg[o*2]+seg[o*2+1];}
void push_down(int l, int r, int o)
{
    if(!tg[0][o] && !tg[1][o])return ;
    int m=(l+r)/2;
    if(tg[0][o])
    {
        seg[o*2]=num[o*2], seg[o*2+1]=num[o*2+1];
        tg[0][o*2]=tg[0][o*2+1]=1;
        tg[1][o*2]=tg[1][o*2+1]=0;
        tg[0][o]=0;
    }
    else
    {
        seg[o*2]=0, seg[o*2+1]=0;
        tg[0][o*2]=tg[0][o*2+1]=0;
        tg[1][o*2]=tg[1][o*2+1]=1;
        tg[1][o]=0;
    }
}
void build(int l, int r, int o)
{
    if(l==r){seg[o]=v[l]-v[l-1], num[o]=seg[o];return ;}
    int m=(l+r)/2;
    build(l, m, o*2), build(m+1, r, o*2+1);
    push_up(o);
    num[o]=num[o*2]+num[o*2+1];
}
void update(int l, int r, int o, int ql, int qr, int k)
{
    if(ql<=l && r<=qr)
    {
        if(k==0)
            seg[o]=num[o], tg[0][o]=1, tg[1][o]=0;
        else 
            seg[o]=0, tg[0][o]=0, tg[1][o]=1;
        return ;
    }
    push_down(l, r, o);
    int m=(l+r)/2;
    if(ql<=m)
        update(l,  m,  o*2,  ql,  qr,  k);
    if(qr>m)
        update(m+1,  r,  o*2+1,  ql,  qr,  k);
    push_up(o);
}
int query(int l, int r, int o, int ql, int qr)
{
    if(ql<=l && r<=qr)
        return seg[o];
    push_down(l, r, o);
    int m=(l+r)/2;
    int ret=0;
    if(ql<=m)
        ret+=query(l, m, o*2, ql, qr);
    if(qr>m)
        ret+=query(m+1, r, o*2+1, ql, qr);
    return ret;
}
main()
{
    scanf("%d%d", &n, &k);
    
    for(int i=1; i<=k; i++)
        q[i].input();
 
    v.push_back(0), v.push_back(n+1);
 
    for(int i = 1; i <= k; i++)
        v.push_back(q[i].l), v.push_back(q[i].r);
 
    sort(v.begin(), v.end());
 
    v.resize( unique(v.begin(), v.end()) - v.begin() );
 
    for(int i=1; i<=k; i++){
        q[i].l = lower_bound(v.begin(),  v.end(),  q[i].l) - v.begin()+1;
        q[i].r = lower_bound(v.begin(),  v.end(),  q[i].r) - v.begin();
    }
    N=v.size()-1;
    build(1, N, 1);
    
    for(int i=1; i<=k; i++){
        update(1, N, 1, q[i].l, q[i].r, q[i].k);
        printf("%d\n", query(1, N, 1, 1, N)-1);
    }
    return 0;
}